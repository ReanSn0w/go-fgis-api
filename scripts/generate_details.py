import json,re,pathlib
info=json.load(open('plans/evidence/registry-api-2026-10-06.json'))
records={}
for e in info['endpoints']:
 if e['path'].endswith('/declarations/{id}'):records['DeclarationDetails']=e['responseSchema']
 if e['path'].endswith('/certificates/{id}'):records['CertificateDetails']=e['responseSchema']
acronyms={'id':'ID','inn':'INN','ogrn':'OGRN','snils':'SNILS','gtin':'GTIN','okpd':'OKPD','tnved':'TNVED','qms':'QMS','eeu':'EEU','ral':'RAL','rpi':'RPI','fgis':'FGIS','ra':'RA','egrul':'EGRUL','opf':'OPF','kpp':'KPP','fias':'FIAS','mrpa':'MRPA','ts':'TS','oon':'OON','un':'UN'}
def name(s):
 parts=re.findall(r'[A-Z]?[a-z]+|[A-Z]+(?=[A-Z]|$)|\d+',s)
 if not parts:parts=[s]
 return ''.join(acronyms.get(p.lower(),p[:1].upper()+p[1:]) for p in parts)
structs={}
def get_type(schema,path):
 types=set(schema['types'])
 nullable='null' in types
 types.discard('null')
 if len(types)!=1:return 'json.RawMessage'
 t=next(iter(types))
 if t=='string':base='string'
 elif t=='integer':base='int64'
 elif t=='number':base='float64'
 elif t=='boolean':base='bool'
 elif t=='object':
  if not schema.get('properties') or any(not name(k)[0].isalpha() for k in schema['properties']):base='map[string]json.RawMessage'
  else:
   base=path
   struct(schema,path)
 elif t=='array':
  item=schema.get('items')
  base='[]'+(get_type(item,path+'Item') if item else 'json.RawMessage')
 else:return 'json.RawMessage'
 if nullable:return '*'+base
 return base

def struct(schema,path):
 if path in structs:return
 lines=[]; names=set()
 for key,child in sorted(schema['properties'].items()):
  field=name(key)
  if field in names:raise RuntimeError((path,key,field))
  names.add(field)
  typ=get_type(child,path+field)
  lines.append(f'\t{field} {typ} `json:"{key}"`')
 if path in ('DeclarationDetails','CertificateDetails'):
  lines.append('\tRaw json.RawMessage `json:"-"`')
 structs[path]=lines
for root,schema in records.items():struct(schema,root)
lines=['// Code generated from the observed 2026-10-06 registry response schemas.','// Values seen only as null or empty arrays remain json.RawMessage.','package fgis','','import "encoding/json"','']
for type_name,fields in sorted(structs.items()):
 lines.extend([f'// {type_name} contains fields observed in public registry detail JSON.',f'type {type_name} struct {{',*fields,'}',''])
for root in records:
 lines.extend([f'// UnmarshalJSON preserves the complete response for fields whose type is unconfirmed.',f'func (d *{root}) UnmarshalJSON(data []byte) error {{',f'\ttype plain {root}','\tvar value plain','\tif err := json.Unmarshal(data, &value); err != nil {','\t\treturn err','\t}',f'\t*d = {root}(value)','\td.Raw = append([]byte(nil), data...)','\treturn nil','}',''])
p=pathlib.Path('details_models.go');p.write_text('\n'.join(lines)+'\n')
print(len(structs),'structs',len(lines),'lines',p.stat().st_size,'bytes')
