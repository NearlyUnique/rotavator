# Profile page

## Features

* See next duty
* manage availability
* request a swap
* link to available swaps days
* link to rota
* edit details

## Non UI features

* email week before next duty
*

## Profile UI

```
Name: <Bob Smith> [📝]
Email: <any@example.com> [📝]
Next Duty: <Saturday 12th June 2023> [Request Swap]
Previous Duty: <Saturday 1st May 2023>
Available on: (Mon ✅, Tue ✅, Wed ✅, Thurs ✅, Fri ❌, Sat ❌, Sun ❌) [📝]
Except : [➕]
  * (Dates) [➖️]
  * (Dates) [➖️]
Conflicts:
  * (Date) [Request Swap]
  * (Date) [Request Swap]

Swaps: <link>
Current Rota: <link>
```

## Other pages
### Swaps
```
# Swaps #
[show all] [offers only] [days-of-week]
|  Date  |  Time  |  Original |  New   | State    |
|--------|--------|-----------|--------|----------|
| <date> | <time> | <name>    | <name> | offered  |
| <date> | <time> | <name>    | <name> | accepted |
```
### Current Rota
```
# Next <7|14> days #
|  Date  |  Time  |  Who   |
|--------|--------|--------|
| <date> | <time> | <name> |
| <date> | <time> | <name> |
```
