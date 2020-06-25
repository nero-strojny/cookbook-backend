import React, { useState } from "react";
import { List, Card, Grid, Icon, Loader } from "semantic-ui-react";
import { deleteRecipe, updateRecipe } from "../serviceCalls";

function RecipeCard({ recipe, refreshRecipesAfterDelete, onEditRecipe }) {

  const {
    recipename: recipeName,
    cooktime: cookTime,
    author,
    _id: recipeId,
    ingredients,
    steps,
    servings,
    rating
  } = recipe;

  const [stepsVisible, setStepsVisible] = useState(false)
  const [ingredientsVisible, setIngredientsVisible] = useState(false)
  const [recipeLoading, setRecipeLoading] = useState(false);
  const [currentRating, setCurrentRating] = useState(rating)

  const onDeleteRecipe = async () => {
    setRecipeLoading(true);
    await deleteRecipe(recipeId);
    refreshRecipesAfterDelete(recipe);
    setRecipeLoading(false);
  }

  const onSelectStar = async (event) => {
    const idValueArray = event.target.id.split('-');
    const submittedRating = parseInt(idValueArray[2]);
    const tempRecipe = { ...recipe, rating: submittedRating }
    await updateRecipe(recipeId, tempRecipe);
    setCurrentRating(submittedRating);
  }

  const createStars = () => {
    const starInputs = [];
    let currentStarNumber = 1;
    for (let i = 0; i < currentRating; i++) {
      starInputs.push(
        <Icon
          name="star"
          color="orange"
          key={"star" + i + recipeId}
          id={`star-${recipeId}-${currentStarNumber++}`}
          onClick={(event) => onSelectStar(event)}
        />)
    }
    for (let i = 0; i < (5 - currentRating); i++) {
      starInputs.push(
        <Icon
          name="star outline"
          color="orange"
          key={"empty star" + i + recipeId}
          id={`star-${recipeId}-${currentStarNumber++}`}
          onClick={(event) => onSelectStar(event)}
        />)
    }
    return (starInputs);
  }

  const createIngredients = () => {
    if (ingredientsVisible) {
      return (
        <>
          <Grid.Row>
            <Grid.Column>
              <h4>
                <Icon
                  style={{ cursor: 'pointer' }}
                  name="minus"
                  color='orange'
                  onClick={() => setIngredientsVisible(false)}
                ></Icon>
                {"\tIngredients"}
              </h4>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column>
              <List bulleted>
                {ingredients.map((ingredient) => (
                  <List.Item key={"ingredient-" + ingredient.name + recipeId}>
                    {ingredient.amount} {ingredient.measurement} of{" "}
                    {ingredient.name}
                  </List.Item>
                ))}
              </List>
            </Grid.Column>
          </Grid.Row>
        </>
      )
    }
    return (
      <Grid.Row>
        <Grid.Column>
          <h4>
            <Icon
              style={{ cursor: 'pointer' }}
              name="plus"
              color='orange'
              onClick={() => setIngredientsVisible(true)}
            ></Icon>
            {"\tIngredients ..."}
          </h4>
        </Grid.Column>
      </Grid.Row>
    );
  }

  const createSteps = () => {
    if (stepsVisible) {
      return (
        <>
          <Grid.Row>
            <Grid.Column>
              <h4>
                <Icon
                  style={{ cursor: 'pointer' }}
                  name="minus"
                  color='orange'
                  onClick={() => setStepsVisible(false)}
                ></Icon>
                {"\tSteps"}
              </h4>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row key={"viewStep" + recipeId}>
            <Grid.Column>
              <List ordered>
                {steps.map((step) => (
                  <List.Item key={"step-text-" + step.number + recipeId}>{step.text}</List.Item>
                ))}
              </List>
            </Grid.Column>
          </Grid.Row>
        </>
      )
    }
    return (
      <Grid.Row>
        <Grid.Column>
          <h4>
            <Icon
              style={{ cursor: 'pointer' }}
              name="plus"
              color='orange'
              onClick={() => setStepsVisible(true)}
            ></Icon>
            {"\tSteps ..."}
          </h4>
        </Grid.Column>
      </Grid.Row>
    );
  }

  return (
    <Card fluid color="orange">
      <Card.Content>
        <Card.Header>
          <Grid>
            <Grid.Row columns="equal">
              <Grid.Column>{recipeName}</Grid.Column>
              <Grid.Column textAlign="right">
                {createStars()}
              </Grid.Column>
            </Grid.Row>
          </Grid>
        </Card.Header>
        {!recipeLoading && <Grid>
          <Grid.Row columns="equal">
            <Grid.Column>
              <Card.Meta>
                <div>By {author}</div>
              </Card.Meta>
            </Grid.Column>
            <Grid.Column >
              <Card.Meta textAlign='right'>
                <div floated='right'>Cook Time: {cookTime} min</div>
                <div>Servings: {servings}</div>
              </Card.Meta>
            </Grid.Column>
          </Grid.Row>
        </Grid>}
        {recipeLoading ?
          <Loader active inline='centered' size='massive' /> :
          <Grid>
            {createIngredients()}
            {createSteps()}
          </Grid>
        }
      </Card.Content>
      <Card.Content extra textAlign="right">
        {!recipeLoading &&
          (<><a onClick={() => onEditRecipe(recipe)}>
            <Icon name="pencil" />
            {`Edit\t`}
          </a>
            <a onClick={() => onDeleteRecipe(recipe)}>
              <Icon name="trash" />
            Delete
          </a></>)
        }
      </Card.Content>
    </Card >
  );
}

export default RecipeCard;
