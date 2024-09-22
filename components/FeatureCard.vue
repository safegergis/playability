<template>
  <Card v-if="game" class="dark p-4 rounded-lg:">
    <h2 class="text-xl font-semibold mb-4">Essential Accessibility Features</h2>
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Colorblind Mode</h3>
        <p class="text-gray-300">
          {{ game.color_blind_support_percentage }}% report full support
        </p>
        <p class="text-gray-300">
          {{ game.color_blind_limited_percentage }}% report limited support
        </p>
        <p class="text-gray-300">
          {{ game.color_blind_no_support_percentage }}% report no support
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Closed Captions</h3>
        <p class="text-gray-300">
          {{ game.closed_captions_support_percentage }}% report full support
        </p>
        <p class="text-gray-300">
          {{ game.closed_captions_limited_percentage }}% report limited support
        </p>
        <p class="text-gray-300">
          {{ game.closed_captions_no_support_percentage }}% report no support
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Controller Support</h3>
        <p class="text-gray-300">
          {{ game.full_controller_support_percentage }}% report full support
        </p>
        <p class="text-gray-300">
          {{ game.controller_limited_support_percentage }}% report limited
          support
        </p>
        <p class="text-gray-300">
          {{ game.controller_no_support_percentage }}% report no support
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Full Controller Remapping</h3>
        <p class="text-gray-300">
          {{ game.controller_remapping_support_percentage }}% report full
          support
        </p>
        <p class="text-gray-300">
          {{ game.controller_remapping_limited_percentage }}% report limited
          support
        </p>
        <p class="text-gray-300">
          {{ game.controller_remapping_no_support_percentage }}% report no
          support
        </p>
      </div>
    </div>
  </Card>
</template>

<script setup></script>

<style></style>
