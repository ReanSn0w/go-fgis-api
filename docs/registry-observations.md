# Наблюдавшиеся справочники ФГИС

При чтении карточек сайт запрашивает `GET /api/v1/rds/common/identifiers` или `GET /api/v1/rss/common/identifiers`, а затем `POST /nsi/api/multi` для выбранных записей классификаторов. Библиотека предоставляет эти методы и пакетное разрешение внутренних ID ТН ВЭД и ОКПД2.

В интерфейсе также наблюдались `POST /nsi/api/oksm/get` (перечень стран) и `POST /nsi/api/status/get` (перечень статусов). Эти полные перечни библиотека пока не запрашивает отдельными публичными методами: для изученных карточек хватает `identifiers` и `multi`. При появлении сценария, требующего перечислить все записи, нужно отдельно проверить форму запроса и пагинацию этих маршрутов.

ID справочной записи и её код — разные значения. Например, `product.identifications[].idTnveds[]` хранит числовой внутренний ID; строковое значение `code` приходит из `tnved` в ответе `multi`. Для страны `oksm.id` имеет строковый тип. Схемы и примеры запросов без данных карточек находятся в [`plans/evidence/registry-api-2026-10-06.json`](../plans/evidence/registry-api-2026-10-06.json).

## Поля с неподтверждённым непустым типом

В исследованных карточках следующие поля были только `null`, пустыми массивами или имели недостаточно примеров для устойчивого типа. Их исходные значения сохранены как `json.RawMessage` либо `[]json.RawMessage`; список отражает выборку, а не запрет на будущие значения:

| Карточка | Поля |
| --- | --- |
| Декларация, корень | `accreditationBody`, `applicantFilials`, `applicationDate`, `applicationNumber`, `applicationSubmissionDate`, `certificationAuthority`, `changes`, `declarationRegInsteadNum`, `declarationReplacedNum`, `experts`, `firstName`, `functions`, `idApplication`, `idApplicationStatus`, `idDeclarationChecker`, `idDeclarationRegInstead`, `idDeclarationReplaced`, `idProductSingleLists`, `idReplacementReason`, `idSigner`, `idSignerEmployee`, `patronymic`, `snils`, `submissionDate`, `surname`, `violationPublishDate` |
| Сертификат, корень | `applicantFilials`, `awaitOperatorCheck`, `changes`, `idCertBasis`, `idDeclarationChecker`, `idFgisRaV1`, `idMrpa`, `idProductSingleLists`, `productInfoOon106`, `productMatchOON106` |
| Продукция декларации | `batchId`, `batchSize`, `identification`, `marking`, `usageCondition`, `usageScope`; у `identifications[]`: `amount`, `article`, `description`, `expiryDate`, `factoryNumber`, `idOkei`, `lifeTime`, `model`, `productionDate`, `sort`, `standards`, `storageTime`, `tradeMark`, `type` |
| Продукция сертификата | `batchId`, `identification`, `marking`; у `identifications[]`: `amount`, `expiryDate`, `factoryNumber`, `idOkpds`, `lifeTime`, `productionDate`, `sort`, `tradeMark`, `type` |

Та же консервативная схема применена к отдельным nullable-полям организаций, документов, протоколов и статусов. Полный JSON каждой карточки доступен в `Raw`. Генератор моделей берёт типы из объединённых схем evidence-файла и позволяет повторить сопоставление после расширения выборки.
