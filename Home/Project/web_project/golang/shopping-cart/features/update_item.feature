Feature: Update Item Quantity

  Scenario Outline: Successfully update the item quantity
    Given I have added a product with ID "<product_id>", name "<product_name>", price "<price>", and stock "<stock>" is available
    And I have added "<quantity>" of product "<product_id>" to the cart 
    When I update the quantity of product "<product_id>" to "<new_quantity>"
    Then the quantity of the product in the cart should be updated to "<new_quantity>"

    Examples:
      | product_id | product_name | price | stock | quantity | new_quantity |
      | 1          | Laptop       | 1000  | 5     | 1        | 3            |
      | 2          | Phone        | 500   | 10    | 2        | 4            |
      | 3          | Desktop      | 1000  | 4     | 1        | 2            |
      | 4          | Headphones   | 50    | 15    | 10       | 7            |

