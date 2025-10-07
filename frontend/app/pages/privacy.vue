<template>
  <main class="min-h-screen dark">
    <!-- Hero Section -->
        <section class="relative dark overflow-hidden bg-gradient-to-b from-background via-background to-neutral-900 py-12"> 
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto text-center space-y-6 animate-fade-in">
          <div class="inline-flex items-center justify-center w-24 h-24 rounded-full bg-primary/10 mb-4">
            <Icon name="lucide:shield-check" class="text-5xl w-14 h-14 text-primary" aria-hidden="true" />
          </div>
          <h1 class="text-4xl lg:text-6xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text">
            Privacy Policy
          </h1>
          <p class="text-xl lg:text-2xl text-muted-foreground leading-relaxed max-w-3xl mx-auto">
            Last Updated: {{ new Date().toLocaleDateString() }}
          </p>
        </div>
      </div>
    </section>

    <!-- Content Section -->
    <section class="py-12 " aria-labelledby="privacy-heading">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto">
          <Card class="dark p-6 lg:p-10 space-y-8">
            <div>

              <h2 class="text-2xl font-bold mb-4 mt-8">1. Introduction</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">Playability ("we", "us", or "our") is committed to protecting your privacy. This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you use our website and services.</p>

              <h2 class="text-2xl font-bold mb-4 mt-8">2. Information We Collect</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">We may collect information about you in a variety of ways. The information we may collect on the Service includes:</p>
              <ul class="text-muted-foreground leading-relaxed mb-4 list-disc pl-6">
                <li><strong>Personal Data:</strong> When you register for an account, we collect personally identifiable information, such as your username and email address.</li>
                <li><strong>User-Generated Content:</strong> We collect the accessibility reports you submit, which includes game information, platform, your score, and any comments you provide.</li>
                <li><strong>Derivative Data:</strong> Information our servers automatically collect when you access the Service, such as your IP address, browser type, and the pages you have viewed.</li>
              </ul>

              <h2 class="text-2xl font-bold mb-4 mt-8">3. How We Use Your Information</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">Having accurate information about you permits us to provide you with a smooth, efficient, and customized experience. Specifically, we may use information collected about you via the Service to:</p>
              <ul class="text-muted-foreground leading-relaxed mb-4 list-disc pl-6">
                <li>Create and manage your account.</li>
                <li>Display your username in connection with your accessibility reports.</li>
                <li>Aggregate anonymous statistical data and analysis for use internally or with third parties.</li>
                <li>Email you regarding your account or verification.</li>
              </ul>

              <h2 class="text-2xl font-bold mb-4 mt-8">4. Data Security</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">We use administrative, technical, and physical security measures to help protect your personal information. While we have taken reasonable steps to secure the personal information you provide to us, please be aware that despite our efforts, no security measures are perfect or impenetrable, and no method of data transmission can be guaranteed against any interception or other type of misuse.</p>

              <h2 class="text-2xl font-bold mb-4 mt-8">5. Cookies</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">We use cookies to manage user sessions for logged-in users. A cookie is a small piece of data stored on your computer. We use a session cookie for authentication purposes. This cookie is essential for the site to function correctly for logged-in users.</p>

              <h2 class="text-2xl font-bold mb-4 mt-8">6. Your Rights</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">You have the right to access, correct, or delete your personal data. You may review or change the information in your account or terminate your account at any time by contacting us.</p>

              <h2 class="text-2xl font-bold mb-4 mt-8">7. Changes to This Privacy Policy</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">We may update this Privacy Policy from time to time. We will notify you of any changes by posting the new Privacy Policy on this page. You are advised to review this Privacy Policy periodically for any changes.</p>
            </div>
          </Card>
        </div>
      </div>
    </section>
  </main>
</template>

<script lang="ts" setup>
useHead({
  title: 'Privacy Policy - Playability',
  meta: [
    {
      name: 'description',
      content: 'Read the Privacy Policy for the Playability website and services.'
    }
  ]
});
</script>
