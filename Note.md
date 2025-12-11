# Plan for Implementing Otel Tracing in Service Modules

## 1. Identify Target locations

Based on the review of the `account_block` module, error handling and important logical branches are the primary candidates for tracing. The goal is to add tracing to functions that interact with external services (like databases, gRPC calls) or perform critical business logic where errors can occur.

Specifically, I will look for:
- Functions returning an `error`.
- Code blocks where errors are logged.
- Calls to other services.

## 2. Tracing Implementation Strategy

The existing tracing implementation in `account_block.go` uses a utility function `local_util.TraceLogger`. I will follow this established pattern.

For each target location, I will implement tracing using the following approach:

-   **Span Creation:** At the beginning of a function, a new `span` will be created by calling `local_util.TraceLogger(ctx, ...)`. The arguments to this function are `(ctx, "layer", "functionName", "serviceGroup", "operation")`.

-   **Span Completion:** The span will be ended at the end of the function using `defer span.End()`.

-   **Event and Error Recording:** When an event of interest or an error occurs, it will be recorded on the span using `span.AddEvent("event name", trace.WithAttributes(...))`. Errors will be included as attributes in these events.

## Example from `account_block`

```go
// This is the pattern I will replicate:
func (s *accountBlockService) GetRegionById(ctx context.Context, id string) (*model.AccountBlock, error) {
    ctx, span := local_util.TraceLogger(ctx, "service", "GetRegionById", "BlockAccount", "GetRegionById")
    defer span.End()

    region, err := s.repo.GetRegionById(ctx, id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            span.AddEvent("Region not found", trace.WithAttributes(attribute.String("id", id), attribute.String("error", err.Error())))
            return nil, errors.New(localization.ErrorRegionNotFound.Code)
        }
        span.AddEvent("Failed to fetch region", trace.WithAttributes(attribute.String("id", id), attribute.String("error", err.Error())))
        return nil, err
    }
    return region, nil
}
```

## Progress

- Tracing implemented for `bps_user_service.go` functions: `Authorize`, `FetchUserByUserCode`, `GetAllBPSUsers`, and `UpdateBpsUser`.
- Tracing implemented for `branch_service.go` function: `Authorize`.
- Tracing implemented for `budget_category_service.go` functions: `Authorize`, `CreateBudgetCategory`, `FetchBudgetCategory`, `FetchBudgetCategoryByID`, `UpdateBudgetCategory`, `DeleteBudgetCategory`, and `EnableOrDisableBudgetCategory`.
- Tracing implemented for `bulk_service.go` functions: `Authorize`, `GetAllBulkServices`, `CheckServiceIsEnabledOrDisabled`, `EnableBulkService`, and `DisableBulkService`.
- Tracing implemented for `cps_action` module:
    - `action_service.go` functions: `CreateCPSAction`, `ApproveCPSAction`, `RejectCPSAction`, `GetCPSActionsByDepartment`, `GetCPSActionByID`, `GetCPSActionByUniqueID`, `GetCPSActionByActionCode`, `RollBack`, `GetActionCountsByDepartemnt`.
    - `dispatcher.go` function: `Authorize`.
- Tracing implemented for `cps_action_role_service.go` functions: `FindAllWithPagination`, `GetByActionCode`, `Create`, `Update`, `Enable`, `Disable`, `Authorize`, `validateUniqueIDs`, `validateUniqueIDsInGroups`, `syncIndices`, and `generateIndices`.
