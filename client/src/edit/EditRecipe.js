import React, { useState } from "react";
import { Form, Divider, Grid, Button, Card } from "semantic-ui-react";
import Steps from "./Steps";
import Ingredients from "./Ingredients";
import { createRecipe } from "../serviceCalls";

function EditRecipe({ onBackToMyRecipes, onSuccessfulCreate }) {
  const [recipeName, setRecipeName] = useState("");
  const [author, setAuthor] = useState("");
  const [steps, setSteps] = useState([""]);
  const [ingredients, setIngredients] = useState([{}]);
  const [cookTime, setCookTime] = useState(0);
  const [prepTime, setPrepTime] = useState(0);
  const [servings, setServings] = useState(0);
  const [isLoading, setIsLoading] = useState(false);

  const submitRecipe = async () => {
    const submittedSteps = [];
    for (let i = 0; i < steps.length; i++) {
      if (steps[i] !== "") {
        submittedSteps.push({ number: i + 1, text: steps[i] });
      }
    }
    const submittedIngredients = [];
    for (let i = 0; i < ingredients.length; i++) {
      if (
        ingredients[i].name !== "" ||
        ingredients[i].measurement !== "" ||
        ingredients[i].amount !== ""
      ) {
        submittedIngredients.push(ingredients[i]);
      }
    }
    setIsLoading(true);
    await createRecipe({
      recipeName,
      author,
      cookTime,
      prepTime,
      servings,
      ingredients: submittedIngredients,
      steps: submittedSteps,
    });
    setIsLoading(false);
    onSuccessfulCreate(recipeName);
  };

  return (
    <Grid padded>
      <Grid.Row columns="equal">
        <Grid.Column>
          <h1> Edit Recipe</h1>
        </Grid.Column>
        <Grid.Column textAlign="right">
          <Button basic color="orange" onClick={() => onBackToMyRecipes()}>
            Back to My Recipes
          </Button>
        </Grid.Column>
      </Grid.Row>
      <Grid.Row>
        <Grid.Column>
          <Card fluid>
            <Card.Content>
              <Form>
                <Grid>
                  <Grid.Row columns="equal">
                    <Grid.Column>
                      <h3>Basics</h3>
                    </Grid.Column>
                  </Grid.Row>
                  <Grid.Row columns="equal">
                    <Grid.Column>
                      <Form.Field>
                        <label>Recipe Name</label>
                        <input
                          placeholder="Recipe Name"
                          defaultValue={recipeName}
                          onChange={(event) =>
                            setRecipeName(event.target.value)
                          }
                        />
                      </Form.Field>
                    </Grid.Column>
                  </Grid.Row>
                  <Grid.Row columns="equal">
                    <Grid.Column>
                      <Form.Field>
                        <label>Author</label>
                        <input
                          placeholder="Author"
                          defaultValue={author}
                          onChange={(event) => setAuthor(event.target.value)}
                        />
                      </Form.Field>
                    </Grid.Column>
                  </Grid.Row>
                  <Grid.Row columns="equal">
                    <Grid.Column>
                      <Form.Field>
                        <label>Prep Time (min)</label>
                        <input
                          type="number"
                          placeholder={0}
                          defaultValue={prepTime}
                          onChange={(event) =>
                            setPrepTime(parseInt(event.target.value))
                          }
                        />
                      </Form.Field>
                    </Grid.Column>
                    <Grid.Column>
                      <Form.Field>
                        <label>Cook Time (min)</label>
                        <input
                          type="number"
                          placeholder={0}
                          defaultValue={cookTime}
                          onChange={(event) =>
                            setCookTime(parseInt(event.target.value))
                          }
                        />
                      </Form.Field>
                    </Grid.Column>
                    <Grid.Column>
                      <Form.Field>
                        <label>Servings</label>
                        <input
                          type="number"
                          placeholder={0}
                          defaultValue={servings}
                          onChange={(event) =>
                            setServings(parseInt(event.target.value))
                          }
                        />
                      </Form.Field>
                    </Grid.Column>
                  </Grid.Row>
                </Grid>
              </Form>
              <Divider />
              <Ingredients
                currentIngredients={ingredients}
                setCurrentIngredients={setIngredients}
              />
              <Divider />
              <Steps currentSteps={steps} setCurrentSteps={setSteps} />
            </Card.Content>
            <Card.Content extra>
              {isLoading ? (
                <Form.Button loading></Form.Button>
              ) : (
                <Form.Button onClick={() => submitRecipe()}>Submit</Form.Button>
              )}
            </Card.Content>
          </Card>
        </Grid.Column>
      </Grid.Row>
    </Grid>
  );
}
export default EditRecipe;
