<script>
  import { Route, Router } from "svelte-routing";
  import HomePage from "./pages/HomePage.svelte";

  export let url = "";
</script>

<Router {url}>
  <Route path="/">
    <HomePage />
  </Route>
</Router>
