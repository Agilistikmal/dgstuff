<script>
  import { navigate } from "svelte-routing";
  import { AuthApi } from "../../lib/api/auth";
  import { onMount } from "svelte";

  let user = $state(null);
  onMount(async () => {
    try {
      user = await AuthApi.me();
      if (user?.role !== "admin") {
        navigate("/auth");
      }
    } catch (err) {
      console.error(err);
      navigate("/auth");
    }
  });
</script>

<div>
  <h1>Admin Page</h1>
</div>
