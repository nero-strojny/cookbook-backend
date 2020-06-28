import React from "react";
import { Form, TextArea, Grid, Button, Icon } from "semantic-ui-react";

function Steps({ currentSteps, setCurrentSteps, onStepDelete }) {
  function changeStepValue(indexSelected, valueInput) {
    const tempSteps = [...currentSteps];
    tempSteps[indexSelected] = valueInput;
    setCurrentSteps(tempSteps);
  }

  function createSteps() {
    const stepInputs = [
      <Grid.Row columns="equal" key={"stepsText0"}>
          <Grid.Column width={14}>
            <Form.Field>
              <label>{"Step 1:"}</label>
              <TextArea
                placeholder="Describe step here..."
                value={currentSteps.length > 0 ? currentSteps[0] : "" }
                onChange={(event) => changeStepValue(0, event.target.value)}
              />
            </Form.Field>
          </Grid.Column>
        </Grid.Row>
    ];
    for (let i = 1; i < currentSteps.length; i++) {
      stepInputs.push(
        <Grid.Row columns="equal" key={"stepsText" + i}>
          <Grid.Column width={14}>
            <Form.Field>
              <label>{`Step ${i+1}:`}</label>
              <TextArea
                placeholder="Describe step here..."
                value={currentSteps[i]}
                onChange={(event) => changeStepValue(i, event.target.value)}
              />
            </Form.Field>
          </Grid.Column>
          <Grid.Column width={2} textAlign="center" verticalAlign="middle">
            <Icon
              name="minus circle"
              size='big'
              color='grey'
              style={{ cursor: 'pointer' }}
              onClick={() => onStepDelete(i)}
            />
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
