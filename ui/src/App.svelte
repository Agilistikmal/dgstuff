<script>
  import { Route, Router } from "svelte-routing";
  import HomePage from "./pages/HomePage.svelte";
  import StuffPage from "./pages/StuffPage.svelte";
  import { AppInfoApi } from "./lib/api/appinfo";
  import { onMount } from "svelte";
  import TransactionPage from "./pages/TransactionPage.svelte";
  import AdminPage from "./pages/admin/AdminPage.svelte";
  import AuthPage from "./pages/auth/AuthPage.svelte";

  let appInfo = $state(null);
  onMount(async () => {
    appInfo = await AppInfoApi.get();
  });

  let { url } = $props();
</script>

<svelte:head>
  <title>dgstuff</title>
  <meta
    name="description"
    content="dgstuff is a platform for selling and buying digital products."
  />
  <meta name="keywords" content="dgstuff, digital products, selling, buying" />
  <meta name="author" content="dgstuff" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta name="robots" content="index, follow" />
  <meta name="googlebot" content="index, follow" />
  <link rel="icon" href={appInfo?.logo_url} />
</svelte:head>

<Router {url}>
  <!-- Public routes -->
  <Route path="/stuff/:slug" let:params>
    <StuffPage slug={params.slug} />
  </Route>
  <Route path="/trx/:id" let:params>
    <TransactionPage id={params.id} />
  </Route>
  <Route path="/">
    <HomePage />
  </Route>

  <!-- Admin routes -->
  <Route path="/admin">
    <AdminPage />
  </Route>

  <!-- Auth routes -->
  <Route path="/auth">
    <AuthPage />
  </Route>
</Router>
