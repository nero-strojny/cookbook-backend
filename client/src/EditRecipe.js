import React from "react";
import { Form, Input, Divider, TextArea, Grid, Button, Icon, Card, List } from 'semantic-ui-react'

function EditRecipe({
    onBackToRecipes
}) {
  return (
    <Grid padded>
        <Grid.Row columns='equal'>
            <Grid.Column>
                <h1> Edit Recipe</h1>
            </Grid.Column>
            <Grid.Column textAlign='right'>
                <Button basic color='orange' onClick={onBackToRecipes}> Back to My Recipes</Button>
            </Grid.Column>
        </Grid.Row>
        <Grid.Row>
            <Grid.Column>
                <Card fluid>
                    <Card.Content>
                        <Form>
                            <Grid>
                            <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                        <h3>Basics</h3>
                                    </Grid.Column>
                                </Grid.Row>
                                <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                        <Form.Field>
                                            <label>Recipe Name</label>
                                            <input placeholder='Recipe Name' />
                                        </Form.Field>
                                    </Grid.Column>
                                </Grid.Row>
                                <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                        <Form.Field>
                                            <label>Author</label>
                                            <input placeholder='Author' />
                                        </Form.Field>
                                    </Grid.Column>
                                </Grid.Row>
                                <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                        <Form.Input label='Prep Time' />
                                    </Grid.Column>
                                    <Grid.Column>
                                        <Form.Input label='Cook Time' />
                                    </Grid.Column>
                                    <Grid.Column>
                                        <Form.Input label='Servings' />
                                    </Grid.Column>
                                </Grid.Row>
                            </Grid>
                            <Divider />
                            <Grid>
                                <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                        <h3>Ingredients</h3>
                                    </Grid.Column>
                                </Grid.Row>
                                <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                        <Form.Input label='Amount' />
                                    </Grid.Column>
                                    <Grid.Column>
                                        <Form.Input label='Measurement' />
                                    </Grid.Column>
                                    <Grid.Column>
                                        <Form.Input label='Name' />
                                    </Grid.Column>
                                </Grid.Row>
                                <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                        <Form.Input />
                                    </Grid.Column>
                                    <Grid.Column>
                                        <Form.Input  />
                                    </Grid.Column>
                                    <Grid.Column>
                                        <Form.Input  />
                                    </Grid.Column>
                                </Grid.Row>
                                <Grid.Row>
                                    <Grid.Column>
                                        <Button>
                                            <Icon name='plus' /> Add Ingredient
                                        </Button>
                                    </Grid.Column>
                                </Grid.Row>
                            </Grid>
                            <Divider />
                            <Grid>
                                <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                        <h3>Steps</h3>
                                    </Grid.Column>
                                </Grid.Row>
                                <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                    <Form.Field
                                        control={TextArea}
                                        label='Step 1'
                                        placeholder='Describe step here...'
                                    />
                                    </Grid.Column>
                                </Grid.Row>
                                <Grid.Row columns ='equal'>
                                    <Grid.Column>
                                    <Form.Field
                                        control={TextArea}
                                        label='Step 2'
                                        placeholder='Describe step here...'
                                    />
                                    </Grid.Column>
                                </Grid.Row>
                                <Grid.Row>
                                    <Grid.Column>
                                        <Button>
                                            <Icon name='plus' /> Add Step
                                        </Button>
                                    </Grid.Column>
                                </Grid.Row>
                            </Grid>
                        </Form>
                    </Card.Content>
                    <Card.Content extra>
                        <Form.Button>Submit</Form.Button>
                    </Card.Content>
                </Card>
            </Grid.Column>
        </Grid.Row>
    </Grid>
  );
}
export default EditRecipe;
