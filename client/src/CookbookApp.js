import React, { useState, useEffect } from "react";
import { Header, Segment, Icon, Container, Message, Transition } from "semantic-ui-react";
import ViewRecipes from "./view/ViewRecipes";
import EditRecipe from "./edit/EditRecipe";

function CookbookApp() {
  const [showEditPage, setShowEditPage] = useState(false);

  const [creationSuccess, setCreationSuccess] = useState(false);
  const [deletionSuccess, setDeletionSuccess] = useState(false);
  const [editSuccess, setEditSuccess] = useState(false);
  const [newRecipeName, setNewRecipeName] = useState("");
  const [deletedRecipeName, setDeletedRecipeName] = useState("");
  const [editedRecipeName, setEditedRecipeName] = useState("");
  const [recipeToEdit, setRecipeToEdit] = useState(defaultRecipe);

  useEffect(() => {
    const timer = setTimeout(() => {
      setEditSuccess(false);
      setDeletionSuccess(false);
      setCreationSuccess(false);
    }, 7000);
    return () => clearTimeout(timer);
  }, [creationSuccess, deletionSuccess, editSuccess]);

  function handleSelectEditRecipe(recipe) {
    setRecipeToEdit(recipe);
    setShowEditPage(true);
  }

  function handleSuccessfulEditRecipe(recipeName) {
    setEditedRecipeName(recipeName);
    setEditSuccess(true);
    setShowEditPage(false);
  }

  function handleCreateRecipe(recipeName) {
    setRecipeToEdit(defaultRecipe);
    setNewRecipeName(recipeName);
    setCreationSuccess(true);
    setShowEditPage(false);
  }

  function handleDeleteRecipe(recipeName) {
    setDeletedRecipeName(recipeName);
    setDeletionSuccess(true);
  }

  return (
    <Container fluid>
      <Segment inverted color="orange">
        <Header as="h1">
          <Icon name="food" />
          Cookbook
        </Header>
      </Segment>
      <Transition visible={creationSuccess} animation="scale" duration={500}>
        <Message
          onDismiss={() => setCreationSuccess(false)}
          header="Recipe created!"
          content={`Recipe, "${newRecipeName}", has been added to the database`}
        />
      </Transition>
      <Transition visible={deletionSuccess} animation="scale" duration={500}>
        <Message
          onDismiss={() => setDeletionSuccess(false)}
          header="Recipe deleted!"
          content={`Recipe, "${deletedRecipeName}", has been removed from the database`}
        />
      </Transition>
      <Transition visible={editSuccess} animation="scale" duration={500}>
        <Message
          onDismiss={() => setEditSuccess(false)}
          header="Recipe edited!"
          content={`Recipe, "${editedRecipeName}", has been edited in the database`}
        />
      </Transition>
      {showEditPage ? (
        <EditRecipe
          onBackToMyRecipes={() => setShowEditPage(false)}
          onSuccessfulCreate={(name) => handleCreateRecipe(name)}
          onSuccessfulEdit={(name) => handleSuccessfulEditRecipe(name)}
          inputtedRecipe={recipeToEdit}
        />
      ) : (
          <ViewRecipes
            onCreateRecipe={() => {
              setShowEditPage(true);
              setRecipeToEdit(defaultRecipe);
            }}
            onSuccessfulDelete={(name) => handleDeleteRecipe(name)}
            onEditRecipe={handleSelectEditRecipe}
          />
        )}
    </Container>
  );
}
export default CookbookApp;

export const defaultRecipe = {
  steps: [""],
  ingredients: [{}],
  recipename: "",
  servings: 0,
  author: "",
  cooktime: 0,
  preptime: 0
};
