Feature: Add items to the shopping cart

  Scenario Outline: Add a product to the cart
    Given a product with ID "<product_id>", name "<product_name>", price "<price>", and stock "<stock>" is available
    When I add "<quantity>" of product "<product_id>" to the cart
    Then the total unique items in the cart should be "<unique_items>"

    Examples:
      | product_id | product_name | price | stock | quantity | unique_items |
      | 1          | Laptop       | 1000  | 5     | 1        | 1            |
      | 2          | Phone        | 500   | 10    | 2        | 2            |
      | 3          | Desktop      | 1000  | 4     | 1        | 3            |
      | 4          | Headphones   | 50    | 15    | 10       | 4            |
