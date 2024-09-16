Feature: Remove items from the shopping cart

  Scenario Outline: Remove an item from the cart
    Given I have added a product with ID "<product_id>", name "<product_name>", price "<price>", and stock "<stock>" is available
    And I have added "<quantity>" of product "<product_id>" to the cart    
    When I checkout the items in the cart
    Then the items should no longer be in the cart

    Examples:
      | product_id | product_name | price | stock | quantity |
      | 1          | Laptop       | 1000  | 5     | 1        | 
      | 2          | Phone        | 500   | 10    | 2        | 
      | 3          | Desktop      | 1000  | 4     | 1        | 
      | 4          | Headphones   | 50    | 15    | 10       | 
