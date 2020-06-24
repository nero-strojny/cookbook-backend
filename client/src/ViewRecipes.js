import React from "react";
import { Input, Grid, Button, Icon, Card, List } from 'semantic-ui-react'

function ViewRecipes({
    onCreateRecipe
}) {
  return (
    <Grid padded>
        <Grid.Row columns='equal'>
            <Grid.Column>
                <Input 
                    fluid
                    placeholder='Search Recipe'
                    icon={<Icon name='search' inverted circular link />}
                />
            </Grid.Column>
            <Grid.Column textAlign='right'>
                <Button color='orange' onClick={onCreateRecipe}>
                    <Icon name='plus' />New Recipe
                </Button>
            </Grid.Column>
        </Grid.Row>
        <Grid.Row>
            <Grid.Column>
            <Card.Group itemsPerRow={1}>
                <Card fluid color='orange' >
                    <Card.Content>
                        <Card.Header>
                            <Grid>
                                <Grid.Row columns ='equal'>
                                <Grid.Column>
                                    Mac and Cheese
                                </Grid.Column>
                                <Grid.Column textAlign='right'>
                                    <Icon name='star' />
                                    <Icon name='star' />
                                    <Icon name='star' />
                                    <Icon name='star' />
                                    <Icon name='star outline' />
                                </Grid.Column>
                                </Grid.Row>
                            </Grid>
                        </Card.Header>
                        <Card.Meta textAlign='right'>
                            <div>Cook Time: 8 min</div>
                            <div>Servings: 2</div>
                        </Card.Meta>
                        <h3><Icon name='minus circle' />Ingredients</h3>
                        <List bulleted>
                            <List.Item>170 g Elbow Macaroni</List.Item>
                            <List.Item>6 oz of Evaporated Milk</List.Item>
                            <List.Item>6 oz of Medium Cheddar Cheese</List.Item>
                        </List>
                        <h3><Icon name='minus circle' />Steps</h3>
                        <List ordered>
                            <List.Item>Place macaroni in a medium saucepan or skillet and add just enough cold water to cover. Add a pinch of salt and bring to a boil over high heat, stirring frequently. Continue to cook, stirring, until water has been almost completely absorbed and macaroni is just shy of al dente, about 6 minutes.</List.Item>
                            <List.Item>Immediately add evaporated milk and bring to a boil. Add cheese. Reduce heat to low and cook, stirring continuously, until cheese is melted and liquid has reduced to a creamy sauce, about 2 minutes longer. Season to taste with more salt and serve immediately.</List.Item>
                        </List>
                    </Card.Content>
                    <Card.Content extra textAlign='right'>
                        <a onClick={onCreateRecipe}><Icon name='pencil' />Edit </a>
                        <a><Icon name='trash' />Delete</a>
                    </Card.Content>
                </Card>
                <Card fluid color='orange' >
                    <Card.Content>
                        <Card.Header>
                            <Grid>
                                <Grid.Row columns ='equal'>
                                <Grid.Column>
                                    Eggplant Rice Bowl
                                </Grid.Column>
                                <Grid.Column textAlign='right'>
                                    <Icon name='star' />
                                    <Icon name='star' />
                                    <Icon name='star' />
                                    <Icon name='star outline' />
                                    <Icon name='star outline' />
                                </Grid.Column>
                                </Grid.Row>
                            </Grid>
                        </Card.Header>
                        <Card.Meta textAlign='right'>
                            <div>Cook Time: 15 min</div>
                            <div>Servings: 4</div>
                        </Card.Meta>
                        <h3><Icon name='plus circle' />Ingredients ...</h3>
                        <h3><Icon name='plus circle' />Steps ...</h3>
                    </Card.Content>
                    <Card.Content extra textAlign='right'>
                        <a onClick={onCreateRecipe}><Icon name='pencil' />Edit </a>
                        <a><Icon name='trash' />Delete</a>
                    </Card.Content>
                </Card>
                <Card fluid color='orange' >
                    <Card.Content>
                        <Card.Header>
                            <Grid>
                                <Grid.Row columns ='equal'>
                                <Grid.Column>
                                    Butternut Squash Ravoli
                                </Grid.Column>
                                <Grid.Column textAlign='right'>
                                    <Icon name='star' />
                                    <Icon name='star' />
                                    <Icon name='star' />
                                    <Icon name='star' />
                                    <Icon name='star' />
                                </Grid.Column>
                                </Grid.Row>
                            </Grid>
                        </Card.Header>
                        <Card.Meta textAlign='right'>
                            <div>Cook Time: 2 hr 5 min</div>
                            <div>Servings: 5</div>
                        </Card.Meta>
                        <h3><Icon name='plus circle' />Ingredients ...</h3>
                        <h3><Icon name='plus circle' />Steps ...</h3>
                    </Card.Content>
                    <Card.Content extra textAlign='right'>
                        <a onClick={onCreateRecipe}><Icon name='pencil' />Edit </a>
                        <a><Icon name='trash' />Delete</a>
                    </Card.Content>
                </Card>
            </Card.Group>
            </Grid.Column>
        </Grid.Row>
    </Grid>
  );
}
export default ViewRecipes;
