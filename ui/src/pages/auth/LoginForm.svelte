<script>
  import { AuthApi } from "../../lib/api/auth";

  let email = $state("");
  let password = $state("");
  let loading = $state(false);
  let error = $state(undefined);

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      loading = true;
      error = undefined;
      const response = await AuthApi.login(email, password);
      console.log(response);
    } catch (err) {
      console.error(err);
      if (err instanceof Error && err.message.includes("Unauthorized")) {
        error = "Invalid email or password";
      } else {
        error = err.message;
      }
    } finally {
      loading = false;
    }
  };
</script>

<form
  method="POST"
  class="border border-gray-200 rounded-lg p-4 border-b-4 border-b-brand space-y-2"
  onsubmit={handleLogin}
>
  <h2 class="text-2xl font-bold mb-4">Login</h2>
  <div>
    <label for="email" class="block">Email</label>
    <input
      type="email"
      id="email"
      name="email"
      bind:value={email}
      placeholder="Email"
      required
      class="w-full border border-gray-300 rounded-lg p-2"
    />
  </div>
  <div>
    <label for="password" class="block">Password</label>
    <input
      type="password"
      id="password"
      name="password"
      bind:value={password}
      placeholder="Password"
      required
      class="w-full border border-gray-300 rounded-lg p-2"
    />
  </div>
  <div class="mt-4">
    {#if error}
      <div class="text-red-500 mb-2">{error}</div>
    {/if}
    <button
      type="submit"
      class="w-full bg-brand text-white rounded-lg p-2"
      disabled={loading}
    >
      {loading ? "Logging in..." : "Login"}
    </button>
  </div>
</form>
