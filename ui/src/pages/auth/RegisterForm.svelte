<script>
  let email = $state("");
  let password = $state("");
</script>

<form
  method="POST"
  class="border border-gray-200 rounded-lg p-4 border-b-4 border-b-brand space-y-2"
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
  <div class="mt-4">
    <button type="submit" class="w-full bg-brand text-white rounded-lg p-2">
      Login
    </button>
  </div>
</form>
