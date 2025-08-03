# 📄 Scraper Service Documentation

Este servicio se encarga de hacer el scraping semanal de los activos de inversión definidos en el sistema principal. El microservicio es independiente del resto de la aplicación y expone endpoints mínimos para lectura de snapshots e ingreso de resultados del scraping.

---

## 📦 Stack Tecnológico

- **Lenguaje**: Go (Golang)
- **Framework HTTP**: Fiber
- **ORM**: GORM
- **Base de datos**: PostgreSQL
- **Dockerizado**: Sí
- **Cache**: Redis

---

## 🗃️ Base de Datos

Schema de la base de datos (prisma -> adaptar a go)

```ts
model User {
  id        String   @id @default(uuid())
  email     String   @unique
  name      String
  password  String
  type      TypeUser @relation(fields: [typeId], references: [id])
  typeId    String
  groups    Group[]
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
}

model TypeUser {
  id          String       @id @default(uuid())
  name        String       @unique
  permissions Permission[]
  users       User[]
}

model Permission {
  id          String   @id @default(uuid())
  name        String   @unique
  description String
  typeUser    TypeUser @relation(fields: [typeUserId], references: [id])
  typeUserId  String
}

model Group {
  id        String         @id @default(uuid())
  name      String
  user      User           @relation(fields: [userId], references: [id])
  userId    String
  type      TypeInvestment @relation(fields: [typeId], references: [id])
  typeId    String
  holdings  Holding[]
  createdAt DateTime       @default(now())
  updatedAt DateTime       @updatedAt
}

model TypeInvestment {
  id           String  @id @default(uuid())
  name         String // Ej: "Cedears", "Criptomonedas", "Acciones"
  scrappingUrl String
  currency     String // Ej: "USD", "ARS"
  groups       Group[]
}

model Holding {
  id               String     @id @default(uuid())
  name             String // Ej: "Apple", "BTC"
  code             String // Ej: "AAPL", "BTC"
  group            Group      @relation(fields: [groupId], references: [id])
  groupId          String
  quantity         Float
  lastPrice        Float?
  earnings         Float?
  relativeEarnings Float? // Porcentaje de ganancia o pérdida en relación al snapshot anterior
  snapshots        Snapshot[]
}

model Snapshot {
  id        String   @id @default(uuid())
  price     Float // Precio del holding en el momento del snapshot
  holding   Holding  @relation(fields: [holdingId], references: [id])
  holdingId String
  quantity  Float // Cantidad de holdings al momento del snapshot
  createdAt DateTime @default(now())
}
```

> Requiere la extensión `uuid-ossp` habilitada en PostgreSQL para `uuid_generate_v4()`.

---

## 🔌 Endpoints REST - Comunicación esperada con el servicio principal (Next.js)

### Validar holding

**Desde Main Service a Snapshot Service:**

```http
POST /api/validate
Authorization: (mediandte API KEY)
```

**body **:

```json
{
  name: string;
  code: string;
  quantity: number;
  groupId: number;
}
```

**Response **:

```json
{
  holding: {
    name: string;
    code: string;
    quantity: number;
    group_id: number;
  };
  is_valid: boolean;
}
```

---

## 🧠 Comportamiento esperado del servicio

1. El scrapping service debe tener un cron que todos los domingos a la 1:00 AM debe buscar en la DB todos los holdings de todos los grupos, buscar el scrapping url del type de cada grupo e ir holding por holding haciendo un scrapping en esa url obteniendo el precio actual de cada activo. Al conseguirlo, se crea una snapshot de ese activo, se actualiza el holding con los earnings y last price y se sigue al siguiente holding.
2. Validate holding es un endpoint que da posibilidad al main service de poder validar al momento de la creacion de un grupo, al agregar un holding si este es valido y es encontrado en la scrapping url del type del grupo. Al recibir una peticion se debe buscar en la url del type si ese activo del holding existe y es valido.

## Envs locales

#LOCAL
DATABASE_URL=postgresql://holding_admin:holding_password@localhost:5432/holdingdb
REDIS_URL=redis://localhost:6379
SNAPSHOT_SERVICE_API_KEY=
