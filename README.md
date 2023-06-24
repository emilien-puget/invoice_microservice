# invoice_microservice

### improvement to be done

- use a transaction to modify balance and mark as paid
- store the transaction in a ledger database
- more tests, including end to end complete scenario

### C4C uml diagram

```mermaid
C4Context
    title github.com/emilien-puget/invoice_microservice
    
    Container_Boundary(user, "user") {
    Component(user.GetAllHandler, "user.GetAllHandler", "", "")
    Component(user.Repository, "user.Repository", "", "")
    
    }
    
    
    Container_Boundary(invoice, "invoice") {
    Component(invoice.CreateInvoiceHandler, "invoice.CreateInvoiceHandler", "", "")
    Component(invoice.Repository, "invoice.Repository", "", "")
    Component(invoice.DoTransactionHandler, "invoice.DoTransactionHandler", "", "")
    
    }
    Rel(user.GetAllHandler, "user.Repository", "GetAll")
    Rel(invoice.CreateInvoiceHandler, "invoice.Repository", "invoice.Repository")
    Rel(invoice.CreateInvoiceHandler, "user.Repository", "GetById")
    Rel(invoice.DoTransactionHandler, "invoice.Repository", "GetByID")
    Rel(invoice.DoTransactionHandler, "invoice.Repository", "MarkAsPaid")
    Rel(invoice.DoTransactionHandler, "user.Repository", "ModifyBalance")
    Component(database_sql.DB, "database_sql.DB", "", "", $tags="external")
    Rel(user.Repository, "database_sql.DB", "database/sql.DB")
    Rel(invoice.Repository, "database_sql.DB", "database/sql.DB")
    Component(github.com_go-playground_validator_v10.Validate, "github.com_go-playground_validator_v10.Validate", "", "", $tags="external")
    Rel(invoice.CreateInvoiceHandler, "github.com_go-playground_validator_v10.Validate", "github.com/go-playground/validator/v10.Validate")
    Rel(invoice.DoTransactionHandler, "github.com_go-playground_validator_v10.Validate", "github.com/go-playground/validator/v10.Validate")

```
