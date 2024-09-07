Feature: Add items to the shopping cart

  Scenario Outline: Add an item to the cart
    When I add an item with ID <itemID>, name "<itemName>", price <price>, and quantity <quantity>
    Then the item should be added successfully
    And the total number of unique items should be <uniqueItemCount>

    Examples:
      | itemID | itemName     | price  | quantity | uniqueItemCount |
      | 1      | "Mobile"     | 150000 | 1        | 1               |
      | 2      | "Desktop     | 30000  | 2        | 2               |
      | 3      | "Mouse"      | 1700   | 1        | 3               |
      | 4      | "Keyboard"   | 5000   | 1        | 4               |