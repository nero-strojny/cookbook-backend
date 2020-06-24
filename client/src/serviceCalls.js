import axios from "axios";

let endpoint = "http://localhost:8080";

export const createRecipe = async (recipe) => {
  console.log(recipe);
  const response = await axios.post(endpoint + "/api/recipe", recipe, {
    headers: {
      "Content-Type": "application/json",
    },
  });
  console.log(response);
};
