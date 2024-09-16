Feature: Apply Discounts to Cart

  Scenario Outline: Apply a discount to the item price
    Given I have added a product with ID "<product_id>", name "<product_name>", price "<price>", and stock "<stock>" is available
    And I have added "<quantity>" of product "<product_id>" to the cart 
    When I apply a discount of <discount_value> with "<discount_type>" discount type to the cart
    Then the expected price should be <expected_price>

    Examples:
      | product_id | product_name | price | stock | quantity | discount_value |discount_type  |expected_price|
      | 1          | Laptop       | 1000  | 5     | 1        | 10             |percentage     |900           |
      | 2          | Phone        | 500   | 10    | 2        | 20             |flat_rate      |480           |
      | 3          | Desktop      | 1000  | 4     | 1        | 10             |flat_rate      |990           |
      | 4          | Headphones   | 50    | 15    | 10       |  2             |percentage     |49            |

