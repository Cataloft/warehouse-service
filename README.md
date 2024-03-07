<h3 align="left">Сервис реализован согласно спецификации JSON:API</h3>



<h3 align="left">Запуск сервиса</h3>

```sh
$ git clone https://github.com/Cataloft/warehouse-service
$ make up
```

<h3 align="left">Доступные ручки</h3>

* **PATCH-`localhost:1234/goods`:** Резервирование или освобождение резерва товара на складе (резервирование: amount>0, освобождение резерва: amount<0)
* **GET-`localhost:1234/warehouse/id`:** Отдаёт данные склада, включая кол-во оставшихся товаров

<h3 align="left">Публичная коллекция Postman</h3>

[Ссылка](https://www.postman.com/descent-module-geoscientist-50761181/workspace/warehouse/collection/29621690-87502690-c891-46e4-afca-1bd37d504e90?action=share&creator=29621690)
