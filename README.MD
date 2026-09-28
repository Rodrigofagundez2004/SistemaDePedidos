## Modelo de datos
  Servicio DE PRODUCTS (Base de datos independiente)
![Diagrama ER](1.1.JPG)
  Servicio de USERS (base de datos independiente)
![Diagrama ER](1.2.JPG)
  Servicio de ORDERS (base de datos independiente)
![Diagrama ER](1.3.JPG)
  Servicio de PAYMENTS (base de datos independiente)
![Diagrama ER](1.4.JPG)

El modelo está separado por microservicio. Las relaciones entre servicios
son referencias lógicas (por ID), no foreign keys físicas.