import React, { useEffect, useState } from "react";
import { Input, Grid, Button, Icon, Card, List } from "semantic-ui-react";
import axios from "axios";

function ViewRecipes({ onCreateRecipe }) {
  const [recipeState, setRecipeState] = useState({});

  const [recipes, setRecipes] = useState([]);

  function toggleIngredients(id) {
      console.log("HERE: " + id)
      console.log(recipeState)
    const tempRecipes = { ...recipeState };
    tempRecipes[id].ingredientsVisible = !tempRecipes[id].ingredientsVisible;
    setRecipeState(tempRecipes);
  }

  function toggleSteps(id) {
    const tempRecipes = { ...recipeState };
    tempRecipes[id].stepsVisible = !tempRecipes[id].stepsVisible;
    setRecipeState(tempRecipes);
  }

  useEffect(() => {
    (async () => {
      let res = await axios.get("http://localhost:8080/api/recipes");
      const tempMap = {};
      res.data.forEach(r => {
        tempMap[r._id] = { ingredientsVisible: true, stepsVisible: true };
      });
      console.log("Jake Was Here")

      setRecipeState(tempMap);
      setRecipes(res.data);
    })();
  }, []);
  return (
    <Grid padded>
      <Grid.Row columns="equal">
        <Grid.Column>
          <Input
            fluid
            placeholder="Search Recipe"
            icon={<Icon name="search" inverted circular link />}
          />
        </Grid.Column>
        <Grid.Column textAlign="right">
          <Button color="orange" onClick={() => onCreateRecipe()}>
            <Icon name="plus" />
            New Recipe
          </Button>
        </Grid.Column>
      </Grid.Row>
      <Grid.Row>
        <Grid.Column>
          <Card.Group itemsPerRow={1}>
            {recipes.map((r) => (
              <Card fluid color="orange">
                <Card.Content>
                  <Card.Header>
                    <Grid>
                      <Grid.Row columns="equal">
                        <Grid.Column>{r.recipename}</Grid.Column>
                        <Grid.Column textAlign="right">
                          <Icon name="star" />
                          <Icon name="star" />
                          <Icon name="star" />
                          <Icon name="star" />
                          <Icon name="star outline" />
                        </Grid.Column>
                      </Grid.Row>
                    </Grid>
                  </Card.Header>
                  <Card.Meta textAlign="right">
                    <div>Cook Time: {r.cooktime} min</div>
                    <div>Servings: {r.servings}</div>
                  </Card.Meta>
                  {recipeState[r._id].ingredientsVisible ? (
                    
                    <Grid.Row><Button
                        icon="minus"
                        basic
                        circular
                        content="Ingredients"
                        labelPosition="left"
                        onClick={() => toggleIngredients(r._id)}
                      ></Button>
                      <List bulleted>
                        {r.ingredients.map((ingredient) => (
                          <List.Item>
                            {ingredient.amount} {ingredient.measurement} of{" "}
                            {ingredient.name}
                          </List.Item>
                        ))}
                      </List>
                    </Grid.Row>
                      
                  ) : (
                    <Grid.Row>
                      <Button
                        icon="plus"
                        basic
                        circular
                        content="Ingredients"
                        labelPosition="left"
                        onClick={() => toggleIngredients(r._id)}
                      ></Button>
                    </Grid.Row>
                  )}

                  {recipeState[r._id].stepsVisible ? (
                    <>
                      <Button
                        icon="minus"
                        basic
                        circular
                        content="Steps"
                        labelPosition="left"
                        onClick={() => toggleSteps(r._id)}
                      ></Button>
                      <List ordered>
                        {r.steps.map((step) => (
                          <List.Item>{step.text}</List.Item>
                        ))}
                      </List>
                    </>
                  ) : (
                    <>
                      <Button
                        icon="plus"
                        basic
                        circular
                        content="Steps"
                        labelPosition="left"
                        onClick={() => toggleSteps(r._id)}
                      ></Button>
                    </>
                  )}
                </Card.Content>
                <Card.Content extra textAlign="right">
                  <a onClick={onCreateRecipe}>
                    <Icon name="pencil" />
                    Edit{" "}
                  </a>
                  <a>
                    <Icon name="trash" />
                    Delete
                  </a>
                </Card.Content>
              </Card>
            ))}
          </Card.Group>
        </Grid.Column>
      </Grid.Row>
    </Grid>
  );
}
export default ViewRecipes;
