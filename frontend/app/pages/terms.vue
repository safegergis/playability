<template>
  <main class="min-h-screen dark">
    <!-- Hero Section -->
    <section class="relative dark overflow-hidden bg-gradient-to-b from-background via-background to-neutral-900 py-12 ">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto text-center space-y-6 animate-fade-in">
          <div class="inline-flex items-center justify-center w-24 h-24 rounded-full bg-primary/10 mb-4">
            <Icon name="lucide:file-text" class="w-14 h-14 text-5xl text-primary" aria-hidden="true" />
          </div>
          <h1 class="text-4xl lg:text-6xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text">
            Terms of Use
          </h1>
          <p class="text-xl lg:text-2xl text-muted-foreground leading-relaxed max-w-3xl mx-auto">
            Last Updated: {{ new Date().toLocaleDateString() }}
          </p>
        </div>
      </div>
    </section>

    <!-- Content Section -->
    <section class="py-12" aria-labelledby="terms-heading">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto">
          <Card class="dark p-6 lg:p-10 space-y-8">
            <div>

              <h2 class="text-2xl font-bold mb-4 mt-8">1. Introduction</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">Welcome to Playability. These Terms of Use ("Terms") govern your access to and use of the Playability website and services (collectively, the "Service"). By accessing or using the Service, you agree to be bound by these Terms.</p>

              <h2 class="text-2xl font-bold mb-4 mt-8">2. User Accounts</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">To submit accessibility reports, you must create an account. You are responsible for safeguarding your account information and for all activities that occur under your account. You agree to notify us immediately of any unauthorized use of your account.</p>

              <h2 class="text-2xl font-bold mb-4 mt-8">3. User Conduct</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">You agree not to use the Service to:</p>
              <ul class="text-muted-foreground leading-relaxed mb-4 list-disc pl-6">
                <li>Post any content that is unlawful, harmful, threatening, abusive, harassing, defamatory, vulgar, obscene, or otherwise objectionable.</li>
                <li>Impersonate any person or entity, or falsely state or otherwise misrepresent your affiliation with a person or entity.</li>
                <li>Post any unsolicited or unauthorized advertising, promotional materials, "junk mail," "spam," or any other form of solicitation.</li>
              </ul>

              <h2 class="text-2xl font-bold mb-4 mt-8">4. User-Generated Content</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">You are solely responsible for the accessibility reports and other content you submit ("User Content"). By submitting User Content, you grant Playability a worldwide, non-exclusive, royalty-free, sublicensable, and transferable license to use, reproduce, distribute, prepare derivative works of, display, and perform the User Content in connection with the Service.</p>
              <p class="text-muted-foreground leading-relaxed mb-4">We have the right, but not the obligation, to monitor and remove any User Content that we, in our sole discretion, find to be in violation of these Terms or otherwise objectionable.</p>

              <h2 class="text-2xl font-bold mb-4 mt-8">5. Disclaimers</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">The Service and all information on it is provided on an "as is" and "as available" basis. Playability makes no warranties, expressed or implied, and hereby disclaims and negates all other warranties including, without limitation, implied warranties or conditions of merchantability, fitness for a particular purpose, or non-infringement of intellectual property or other violation of rights.</p>
              <p class="text-muted-foreground leading-relaxed mb-4">The accessibility information on this site is sourced from third parties and community contributions and is for informational purposes only. We do not guarantee the accuracy of this information.</p>

              <h2 class="text-2xl font-bold mb-4 mt-8">6. Limitation of Liability</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">In no event shall Playability or its suppliers be liable for any damages (including, without limitation, damages for loss of data or profit, or due to business interruption) arising out of the use or inability to use the materials on Playability's website, even if Playability or a Playability authorized representative has been notified orally or in writing of the possibility of such damage.</p>

              <h2 class="text-2xl font-bold mb-4 mt-8">7. Changes to These Terms</h2>
              <p class="text-muted-foreground leading-relaxed mb-4">We may revise these Terms of Use at any time without notice. By using this website you are agreeing to be bound by the then current version of these Terms of Use.</p>
            </div>
          </Card>
        </div>
      </div>
    </section>
  </main>
</template>

<script lang="ts" setup>
useHead({
  title: 'Terms of Use - Playability',
  meta: [
    {
      name: 'description',
      content: 'Read the Terms of Use for using the Playability website and services.'
    }
  ]
});
</script>
