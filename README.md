# gobud
Budget check and generator tool

usage 

gobud budget check data.csv : som + eom by default
gobud budget set category 1000 : set the value of a budget
gobud budget budget list : list all budgets
gobud budget remove category : remove a budget

gobud init : initialize the field value

gobud import ofx
    - add all transaction to a ledger.json file
    - for every unknown transaction ask for a category
    - List all category before asking for category

# test

gobud list-budget
gobud budget list
gobud budget set
gobud budget report  --start --end
gobud budget delete

gobug tx import
gobud tx category 1 alimentation
gobud tx list --start=test --end=end
