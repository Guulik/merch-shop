
# Введение

TODO: настроить передачу секрета

# Архитектура проекта
Заметил, что все используют разные названия для описания одних и тех же слоёв. Ниже напишу, что я буду иметь в виду под каждым слоем:
- API слой - он же Handlers, Controllers. Это слой работы с пользователем. Он принимает http-запросы и 
передает их дальше в сервисный слой
- Сервис - он же Usecase. Это слой бизнес-логики. Здесь проверяются все условия и здесь же мы выбираем, какую ошибку отдать пользователю,
оборачиваем её, сохраняем контекст для [логирования](#логирование-и-пользовательские-ошибки)
и отдаем обратно в api слой. 
- Repository - он же storage, Репо. Это обычный слой работы с базой данных. Здесь пишем SQL запросы, т.к. в нашем случае
postgres. В принципе не проблема "прикрутить" кеширование в Redis или самописное хранилище, но этого нет в требованиях, 
а время ограничено, поэтому я не стал:)

Репо и сервисный слои покрыты тестами и для них описаны интерфейсы, чтобы слой выше мог не зависеть от конкретной реализации
используемого слоя.

# Тесты
Требование по 40%+ покрытия выполнено :) 
Для проверки тестов и процента покрытия можно использовать:
- `make coverage_cli`, чтобы посмотреть в консоли. Здесь же можно увидеть total %
- `make coverage_html`, чтобы посмотреть в браузере. Откроется интерактивная страница, которая покажет, 
какие участки кода покрыты

Процент покрытия сейчас 61% 
Написаны юниты на все сценарии бизнес-логики для сервисного слоя и для слоя работы с базой (репозиторий).
Чтобы адекватно замокать репо слой, пришлось использовать либу pgxmock и заменить в сервисе зависимость pgxpool.Pool 
на самописный интерфейс. Было интересно подгонять интерфейс под мок под фактическую реализацию pgxpool.Pool )
круто ощутил силу интерфейсов.


# Логирование и пользовательские ошибки
За основу я взял логгер с полным контекстом запроса, который увидел у Алексея Мичурина в 
одном из видео AvitoTech на ютубе. 
Теперь очень удобно можно видеть весь контекст: что за пользователь делает запрос, какой предмет хочет купить,
кому хочет переслать монетки, сколько монеток и т.д.
Посмотреть можно [здесь](./internal/util/logger) 

А пользователю я отдаю кастомные обёрнутые ошибки. 
Например, если ошибка - "no rows in result set", то я говорю:
"user not found" или "receiver does not exist" в зависимости от запроса.
А если ошибка внутренняя, то просто "internal server error" без подробностей.

