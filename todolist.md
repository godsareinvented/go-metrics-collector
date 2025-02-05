# todo list

## Список правок на будущее

---

### Временные

1. [config_configurator](internal/server/config/config_configurator.go:#L22) Необходимо валидация отрицательных значений
   флагов?
2. [config_configurator](internal/server/config/config_configurator.go:#L26) Необходимо разделить логику сервера и
   агента,
   т.к. для агента не требуется инициализация хранилища.
3. [config_configurator](internal/server/config/config_configurator.go:#L26) Значение GzipMinContentLength должно быть
   1400.
   Для соответствия инкременту 8 заменено на 0.
4. [update_metric_json](internal/agent/client/request/update_metric_json.go:#L16) Нужно пересмотреть выплёвывание
   ошибок.
5. [on_start_callback](internal/server/server1/callback/on_start_callback.go:#L20) В будущем, необходимо перейти
   на более надёжную схему менеджера джоб.
6. [server](internal/server/server.go:#L47) Реакция на завершение контекста сервера в хендлере DbPing прописано верно.
   Но остановка сервера в текущем варианте останавливает хендлеры моментально, не давая им возможности корректно
   обработать завершение родительного контекста (точно?).
   Нужно перейти на схему с грациозным завершением сервера (?..). Надо исследовать этот момент глубже.
7. [metric_manager](internal/server/service/metric/metric_manager.go:#L48) Упростить и улучшить именование методов и
   структур.
8. [metric_manager](internal/server/service/metric/metric_manager.go) Общение между горутинами необходимо переписать на
   каналы
    вместо использования переменной пакета.

---

### Валидация

9. [custom_func](internal/general/validation/custom_func): Добавить пользовательские констреинты ```Integer``` и
   ```Float```,
   поддерживающие проверку, что число во входной строке не больше ```math.MaxInt64``` и ```math.MaxFloat64```
   (*плюс поддержка отрицательных значений через регулярку*)\
   \
   Алгоритм следующий: число с количеством разрядов больше чем у максимального значения фильтруется регуляркой.
   Для проверки при точном совпадении количества разрядов пройтись от старшего разряда к младшему,
   сравнивая поразрядно с максимальным значением типа. Если значение хоть какого-то разряда входного числа больше
   значения разряда максимального значения, считать валидацию неуспешной

---

### Хранилища

10. [postgres/storage](internal/server/storage/postgres/storage.go:#L155) Следует ли добавить в будущем контекст с
    дедлайном?

---

### Кодстайл

11. [metric_manager_test](internal/server/buisness_logic/manager/metric_manager_test.go:L72) Заменить все условия вида
   ```if metric.Type == dictionary.GaugeMetricType { // ... } else``` на switch с обработкой ситуации в ```default```
   -секции,
   что метрика имеет некорректный тип
12. [update_metric_batch_handler](internal/server/handler/update_metric_batch_handler.go:#L29) Нужно пересмотреть
    инициализацию объектов
    (убрать постоянное выделение памяти в обработчиках)

---

### Проблемы сервиса

13. Теоретически, для значений метрик с типом ```Counter``` может произойти переполнение переменной

---