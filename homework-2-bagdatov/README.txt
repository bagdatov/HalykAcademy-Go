Необходимо реализовать выполнение последовательных перемещений в ЛОКАЦИЯХ объекта СОТРДУНИК

1) Внутри файла employeeActions находятся строки в формате CSV (разделитель - запятая).
    1. Первая колонка - тип сотрудника. Все последующие это пайплайн его передвижения.
    2. Нужно прочесть файл и записать данные в удобную оболочку для дальнейших действий.

2) Так же вам даны несколько готовых структур: 2 типа СОТРУДНИКОВ и 3 типа ЛОКАЦИЙ.
    1. Им требуется реализация интерфейсов.
        - Для hr и itSecure - Employee
        - Для office, workArea, servers - Location
    2. Для создания этих объектов можно использовать только готовые фабричные методы - NewEmployee и NewLocation.

3) Интерфейс Location имеет следующие ограничения. Перемещаться можно только в соседние локации
    (это уже реализовано в new* методах).
    Если у объекта Employee локации еще нет, то первой может быть только office.

4) У сотрудников есть ограничения по доступу в помещения. Они также уже указаны в new* методах.


- Ваша задача, если вы за нее возьметесь, мистер Хант:
    1. Прочитать файл, и положить данные "куда-то в удобное место"
    2. Итеративно пройтись по каждому набору данных. У вас есть сотрудник и у него есть определенные действия по перемещению
        с ограничениями.
    3. Нужно выводить в консоль каждое успешно выполненное действие
    4. Если возникла ошибка ее нужно обработать, чтобы было понятно откуда она взялась

Ограничения:
    - Нельзя добавлять новые типы
    - Все методы интерфейсов должны быть использованы
    - Нельзя модифицировать исходный файл employeeActions
    - Нельзя модифицировать фабричные функции и функции генерации объектов
    - Программа не должна падать с паниками
    - Нужно пройтись по всем строкам файла, даже если будут возвращаться ошибки - их нужно обработать и продолжить работу