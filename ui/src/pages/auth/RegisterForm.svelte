<script>
  import { navigate } from "svelte-routing";
  import { AuthApi } from "../../lib/api/auth";

  let email = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let loading = $state(false);
  let error = $state(undefined);

  const handleRegister = async (e) => {
    e.preventDefault();

    if (password !== confirmPassword) {
      error = "Passwords do not match";
      return;
    }

    try {
      loading = true;
      error = undefined;
      const response = await AuthApi.register(email, password, confirmPassword);
      cookieStore.set("auth_token", response.token);
      const getMe = await AuthApi.me();
      if (getMe.role === "admin") {
        navigate("/admin");
      } else {
        navigate("/");
      }
    } catch (err) {
      console.error(err);
      error = err.message;
    } finally {
      loading = false;
    }
  };
</script>

<form
  method="POST"
  class="border border-gray-200 rounded-lg p-4 border-b-4 border-b-brand space-y-2"
  onsubmit={handleRegister}
>
  <h2 class="text-2xl font-bold mb-4">Register</h2>
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
  <div>
    <label for="confirmPassword" class="block">Confirm Password</label>
    <input
      type="password"
      id="confirmPassword"
      name="confirmPassword"
      bind:value={confirmPassword}
      placeholder="Confirm Password"
      required
      class="w-full border border-gray-300 rounded-lg p-2"
    />
  </div>
  <div class="mt-4">
    {#if error}
      <div class="text-red-500 mb-2">{error}</div>
    {/if}
    <button type="submit" class="w-full bg-brand text-white rounded-lg p-2">
      {loading ? "Registering..." : "Register"}
    </button>
  </div>
</form>
