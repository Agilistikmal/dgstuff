<script>
  import { Route, Router } from "svelte-routing";
  import HomePage from "./pages/HomePage.svelte";
  import StuffPage from "./pages/StuffPage.svelte";

  export let url = "";
</script>

<Router {url}>
  <Route path="/">
    <HomePage />
  </Route>
  <Route path="/stuff/:slug" let:params>
    <StuffPage slug={params.slug} />
  </Route>
</Router>
