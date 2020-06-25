import React, { useState } from "react";
import { Header, Segment, Icon, Container } from "semantic-ui-react";
import ViewRecipes from "./ViewRecipes";
import EditRecipe from "./edit/EditRecipe";

function CookbookApp() {
  const [showEditPage, setShowEditPage] = useState(false);

  return (
    <Container fluid>
      <Segment inverted color="orange">
        <Header as="h1">
          <Icon name="food" />
          Cookbook
        </Header>
      </Segment>
      {showEditPage ? (
        <EditRecipe onBackToMyRecipes={() => setShowEditPage(false)} />
      ) : (
        <ViewRecipes onCreateRecipe={() => setShowEditPage(true)} />
      )}
    </Container>
  );
}
export default CookbookApp;
