//
//WHO    = Principal (client_id)
//WHAT   = Action (inventory create)
//WHICH  = Resource (product sku-10) when empty means all
//

permit(
    principal,
    action == Action::"get:product",
    resource
)
when {
    principal.scopes.contains(
        "inventory:read"
    )
};

permit(
    principal,
    action == Action::"add:product",
    resource
)
when {
    principal.scopes.contains(
        "inventory:create"
    )
};

permit(
    principal,
    action == Action::"update:product",
    resource
)
when {
    principal.scopes.contains(
        "inventory:update"
    )
};

permit(
    principal,
    action == Action::"patch:inventory",
    resource
)
when {
    principal.scopes.contains(
        "inventory:update"
    )
};