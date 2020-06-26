import React, { useEffect, useState } from "react";
import { Input, Grid, Button, Icon, Card, Transition } from "semantic-ui-react";
import RecipeCard from "./RecipeCard";
import { getRecipes } from "../serviceCalls";

function ViewRecipes({ onCreateRecipe, onSuccessfulDelete, onEditRecipe }) {
  const [recipes, setRecipes] = useState([]);
  const [shouldRefresh, setShouldRefresh] = useState(true);

  useEffect(() => {
    let isCurrent = true;
    (async () => {
      if (isCurrent) {
        if (shouldRefresh) {
          window.scrollTo(0, 0)
          setRecipes(await getRecipes());
        }
        setShouldRefresh(false);
      }
    })();
    return () => {
      isCurrent = false
    }
  }, [shouldRefresh]);

  function refreshRecipesAfterDelete(recipe) {
    setShouldRefresh(true);
    onSuccessfulDelete(recipe.recipename);
  }

  function onRefreshRecipes() {

  }

  return (
    <Grid padded>
      <Grid.Row columns="equal">
        <Grid.Column>
          <Input
            fluid
            placeholder="Search Recipe"
            icon={<Icon name="search" color='orange' inverted circular link />}
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
            {shouldRefresh ?
              (<p style={{ color: 'grey', cursor: 'pointer' }}>
                <Icon
                  loading
                  color='grey'
                  name="spinner"
                ></Icon>{"\tRefresh Recipes"}
              </p>) :
              (<p style={{ color: 'grey', cursor: 'pointer' }} 
                  onClick={() => setShouldRefresh(true)}>
                  <Icon
                    name="refresh"
                    color='grey'
                  ></Icon>{"\tRefresh Recipes"}
              </p>
             )
            }
        </Grid.Column>
      </Grid.Row>
      <Grid.Row>
        <Grid.Column>
          <Card.Group itemsPerRow={1}>
            <Transition.Group
              duration={1500}
            >
              {recipes.map((r) => (
                <RecipeCard
                  recipe={r}
                  refreshRecipesAfterDelete={refreshRecipesAfterDelete}
                  onEditRecipe={onEditRecipe}
                  key={"recipeCard" + r._id}
                />
              ))}
            </Transition.Group>
          </Card.Group>
        </Grid.Column>
      </Grid.Row>
    </Grid>
  );
}
export default ViewRecipes;
