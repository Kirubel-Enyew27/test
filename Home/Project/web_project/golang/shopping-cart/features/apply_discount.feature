Feature: Apply Discounts to Cart

  Scenario Outline: Apply a discount to the item price
    Given I have added a product with ID "<product_id>", name "<product_name>", price "<price>", and stock "<stock>" is available
    And I have added "<quantity>" of product "<product_id>" to the cart 
    When I apply a discount of <discount_value>% to the cart
    Then the expected price should be <expected_price>

    Examples:
      | product_id | product_name | price | stock | quantity | discount_value |expected_price|
      | 1          | Laptop       | 1000  | 5     | 1        | 10             |900           |
      | 2          | Phone        | 500   | 10    | 2        | 20             |400           |
      | 3          | Desktop      | 1000  | 4     | 1        | 10             |900           |
      | 4          | Headphones   | 50    | 15    | 10       |  2             |49            |

