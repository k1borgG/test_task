**Test_task**

**Instructions for start:**

run:
1) "make up"

App ready for use!

**API**

**1:**
Адресс: "localhost:8080"
Метод: "AddProduct"

Добавляет новый продукт в систему.

**Request**
```
{
    string name = 1;
    string description = 2;
    string brand = 3;
    string model = 4;
    Coordinates coordinates = 5;
    google.protobuf.Timestamp created_at = 6;
    uint32 price = 7;
}
```

**Request example**
```
{
  "product": {
    "name": "iPhone",
    "description": "Latest model",
    "brand": "Apple",
    "model": "iPhone 12",
    "coordinates": { "lat": 37.7749, "lon": 122.4194 },
    "created_at": { "seconds": 1609459200 },
    "price": 999
  }
}
```

**Respone example**
```
{
  "id": "@2fdsadf12345Uafafa2234dasfsaf6",
  "product": {
    "name": "iPhone",
    "description": "Latest model",
    "brand": "Apple",
    "model": "iPhone 12",
    "coordinates": { "lat": 37.7749, "lon": 122.4194 },
    "created_at": { "seconds": 1609459200 },
    "price": 999
  }
}
```
----------------------------------------------------------------------------------------------
**2:**
Адресс: "localhost:8080"
Метод: "GetProduct"

Извлекает продукты по заданным фильтрам.

**Request**
```
{
    optional string name = 1;
    optional string description = 2;
    optional string brand = 3;
    optional string model = 4;
    optional Coordinates coordinates = 5;
    optional google.protobuf.Timestamp created_at = 6;
    optional uint32 price = 7;
    optional string filter_date1 = 8;
    optional string filter_date2 = 9;
    optional string filter_price1 = 10;
    optional string filter_price2 = 11;
}
```

**Request example**
```
{
    "product": {
        "name":"iphone",
        "filter_price1": "520",
        "filter_price2": "1000"
    }
}
```

**Respone example**
```
{
  "product": [
    {
      "name": "iPhone",
      "description": "Old model",
      "brand": "Apple",
      "model": "iPhone X",
      "coordinates": { "lat": 37.7749, "lon": 122.4194 },
      "created_at": { "seconds": 1577854800 },
      "price": 700
    }
  ]
}
```
