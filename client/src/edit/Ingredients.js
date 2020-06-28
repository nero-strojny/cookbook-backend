import React from "react";
import { Form, Grid, Button, Icon } from "semantic-ui-react";

function Ingredients({ currentIngredients, setCurrentIngredients, onIngredientDelete }) {

  function changeIngredientValue(indexSelected, keySelected, valueInput) {
    const tempIngredients = [...currentIngredients];
    tempIngredients[indexSelected][keySelected] = valueInput;
    setCurrentIngredients(tempIngredients);
  }

  function createIngredients() {
    let ingredientInputs = [
      <Grid.Row columns="equal" >
        <Grid.Column width={4}>
          <Form.Field>
            <label>Amount</label>
            <input
              type="number"
              defaultValue={
                currentIngredients.length >= 1
                  ? currentIngredients[0].amount
                  : ""
              }
              onChange={(event) =>
                changeIngredientValue(0, "amount", parseFloat(event.target.value))
              }
            />
          </Form.Field>
        </Grid.Column>
        <Grid.Column width={4}>
          <Form.Field>
            <label>Measurement</label>
            <input
              defaultValue={
                currentIngredients.length >= 1
                  ? currentIngredients[0].measurement
                  : ""
              }
              onChange={(event) =>
                changeIngredientValue(0, "measurement", event.target.value)
              }
            />
          </Form.Field>
        </Grid.Column>
        <Grid.Column width={6}>
          <Form.Field>
            <label>Name</label>
            <input
              defaultValue={
                currentIngredients.length >= 1 ? currentIngredients[0].name : ""
              }
              onChange={(event) =>
                changeIngredientValue(0, "name", event.target.value)
              }
            />
          </Form.Field>
        </Grid.Column>
      </Grid.Row>,
    ];

    for (let i = 1; i < currentIngredients.length; i++) {
      ingredientInputs.push(
        <Grid.Row columns="equal" key={"steps" + i}>
          <Grid.Column width={4}>
            <Form.Field>
              <input
                type="number"
                defaultValue={currentIngredients[i].amount}
                onChange={(event) =>
                  changeIngredientValue(i, "amount", parseFloat(event.target.value))
                }
              />
            </Form.Field>
          </Grid.Column>
          <Grid.Column width={4}>
            <Form.Field>
              <input
                defaultValue={currentIngredients[i].measurement}
                onChange={(event) =>
                  changeIngredientValue(i, "measurement", event.target.value)
                }
              />
            </Form.Field>
          </Grid.Column>
          <Grid.Column width={6}>
            <Form.Field>
              <input
                defaultValue={currentIngredients[i].name}
                onChange={(event) =>
                  changeIngredientValue(i, "name", event.target.value)
                }
              />
            </Form.Field>
          </Grid.Column>
          <Grid.Column width={2} textAlign="center" verticalAlign="middle">
            <Icon 
              name="minus circle"
              size='big'
              color='grey'
              style={{ cursor: 'pointer' }}
              onClick={() => onIngredientDelete(i)}
            />
          </Grid.Column>
        </Grid.Row>
      );
    }
    return ingredientInputs;
  }

  return (
    <Form>
      <Grid>
        <Grid.Row columns="equal">
          <Grid.Column>
            <h3>Ingredients</h3>
          </Grid.Column>
        </Grid.Row>
        {createIngredients()}
        <Grid.Row>
          <Grid.Column>
            <Button
              onClick={() => setCurrentIngredients([...currentIngredients, {}])}
            >
              <Icon name="plus" /> Add Ingredient
            </Button>
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Form>
  );
}
export default Ingredients;
