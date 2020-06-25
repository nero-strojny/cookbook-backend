import React, { useState } from "react";
import { Header, Segment, Icon, Container, Message, Transition } from "semantic-ui-react";
import ViewRecipes from "./ViewRecipes";
import EditRecipe from "./edit/EditRecipe";

function CookbookApp() {
  const [showEditPage, setShowEditPage] = useState(false);

  const [creationSuccess, setCreationSuccess] = useState(false);
  const [newRecipeName, setNewRecipeName] = useState("false");

  function handleDismiss() {
    setCreationSuccess(false);
  }

  function handleCreateRecipe(recipeName) {
    setNewRecipeName(recipeName);
    setCreationSuccess(true);
    setShowEditPage(false);
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
            onDismiss={() => handleDismiss()}
            header="Recipe created!"
            content={`Recipe, "${newRecipeName}", has been added to the database`}
          />
      </Transition>
      {showEditPage ? (
        <EditRecipe
          onBackToMyRecipes={() => setShowEditPage(false)}
          onSuccessfulCreate={(name) => handleCreateRecipe(name)}
        />
      ) : (
        <ViewRecipes onCreateRecipe={() => setShowEditPage(true)} />
      )}
    </Container>
  );
}
export default CookbookApp;
