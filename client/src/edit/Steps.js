import React from "react";
import { Form, TextArea, Grid, Button, Icon } from "semantic-ui-react";

function Steps({ currentSteps, setCurrentSteps }) {
  function changeStepValue(indexSelected, valueInput) {
    const tempSteps = [...currentSteps];
    tempSteps[indexSelected] = valueInput;
    setCurrentSteps(tempSteps);
  }

  function createSteps() {
    const stepInputs = [];
    for (let i = 0; i < currentSteps.length; i++) {
      stepInputs.push(
        <Grid.Row columns="equal" key={"steps" + i}>
          <Grid.Column>
            <Form.Field>
              <label>{"Step " + (i + 1)}</label>
              <TextArea
                placeholder="Describe step here..."
                value={currentSteps[i]}
                onChange={(event) => changeStepValue(i, event.target.value)}
              />
            </Form.Field>
          </Grid.Column>
        </Grid.Row>
      );
    }
    return stepInputs;
  }

  return (
    <Form>
      <Grid>
        <Grid.Row columns="equal">
          <Grid.Column>
            <h3>Steps</h3>
          </Grid.Column>
        </Grid.Row>
        {createSteps()}
        <Grid.Row>
          <Grid.Column>
            <Button onClick={() => setCurrentSteps([...currentSteps, ""])}>
              <Icon name="plus" /> Add Step
            </Button>
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Form>
  );
}

export default Steps;
