## How to guide for Avro DTO's....
In this guide I will try to explain how to generate/update avro schemas and some tricks.
Do not forget that we will run everything in the project root.

## How to add new Schema DTO?
Connect to AWS Dev vpn and run this command.

```bash
go run models/avro-dto/main.go -schema=<Key of your schema> -version=<version as int> -package=<package name !(not mandatory)>
```

## How to create new one to replace? (this is not version update)
Connect to AWS Dev vpn and run this command.

```bash
go run models/avro-dto/main.go -schema=<Key of new schema> -version=<version as int> -package=<package name !(not mandatory)>
```

Then find the references of old one by looking for this line.

```
"github.com/yusufpapurcu/papel/models/avro-dto/<key of old schema>"
```

After being sure everything updated and tests passing, you can remove the old one's directory.

## How to do version update?
Learn the version you want to generate, and be sure it's compatible with the one in `schema-versions.env`

Also learn the name of the exiting package. You can check one of the generated files.

Connect to AWS Dev vpn and run this command.

```bash
go run models/avro-dto/main.go -schema=<Key of new schema> -version=<version as int> -package=<package name>
```

By keeping package name consistent, you don't need to do renames. Just be sure you reflect changes you do in the code!
## How to map nullable field into generated DTO?
So let's look into an example case:

```go
type CustomerContext string

const (
	CustomerContextEcom      CustomerContext = "ECOM"
	CustomerContextTerminal  CustomerContext = "TERMINAL"
	CustomerContextPhone     CustomerContext = "PHONE"
	CustomerContextPayByLink CustomerContext = "PAY_BY_LINK"
)

type CanonicalTransaction struct {
    CustomerContext   *CustomerContext             `json:"customer_context"`
}
```

So in our domain model we have `CustomerContext` which is a nullable enum.
So in order to map this data into schema dto, you will need to use a block like this:

```go
if data.CustomerContext != nil {
    res.Customer_context = avro_dto.NewCustomer_contextUnion()
    res.Customer_context.UnionType = avro_dto.Customer_contextUnionTypeEnumAvro_dto_customer_context
    res.Customer_context.Avro_dto_customer_context, err = avro_dto.NewAvro_dto_customer_contextValue(string(*data.CustomerContext))
    if err != nil {
        return avro_dto.Avro_dto{}, fmt.Errorf("failed to map customer_context field: %w", err)
    }
}
```

1. We need to check if field is filled in.
2. We need to initialize new union for customer context and set union type into `customer_contextUnion`
3. After we can put our enum inside. I used auto-generated `NewAvro_dto_customer_contextValue` function, and checked error value ofc.
4. If you want to keep continue the process, you can just set `CustomerContext` field into nil.
