<template>
  <main class="min-h-screen dark">
    <!-- Hero Section -->
    <section class="relative dark overflow-hidden bg-gradient-to-b from-background via-background to-neutral-900 py-16 lg:py-24">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto text-center space-y-6 animate-fade-in">
          <div class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-primary/10 mb-4">
            <Icon name="lucide:heart" class="w-10 h-10 text-primary animate-pulse" aria-hidden="true" />
          </div>
          <h1 class="text-4xl lg:text-6xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text">
            Support Playability
          </h1>
          <p class="text-xl lg:text-2xl text-muted-foreground leading-relaxed max-w-3xl mx-auto">
            Help us keep Playability free and accessible for everyone in the gaming community
          </p>
        </div>
      </div>
    </section>

    <!-- Why Donate Section -->
    <section class="py-16 lg:py-20 dark" aria-labelledby="why-donate-heading">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto space-y-8">
          <div class="text-center space-y-4 animate-slide-up">
            <h2 id="why-donate-heading" class="text-3xl lg:text-4xl font-bold">Why Your Support Matters</h2>
            <p class="text-lg text-muted-foreground">
              Playability is built and maintained by a single developer with a passion for accessibility
            </p>
          </div>

          <div class="grid md:grid-cols-2 gap-6 animate-slide-up" style="animation-delay: 150ms">
            <!-- Reason 1 -->
            <Card class="dark p-6 lg:p-8 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:border-primary/50">
              <div class="flex items-start gap-4">
                <div class="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0">
                  <Icon name="lucide:server" class="w-6 h-6 text-primary" aria-hidden="true" />
                </div>
                <div>
                  <h3 class="text-xl font-bold mb-2">Server Costs</h3>
                  <p class="text-muted-foreground leading-relaxed">
                    Keep our servers running smoothly to provide fast, reliable access to accessibility information for thousands of games.
                  </p>
                </div>
              </div>
            </Card>

            <!-- Reason 2 -->
            <Card class="dark p-6 lg:p-8 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:border-primary/50">
              <div class="flex items-start gap-4">
                <div class="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0">
                  <Icon name="lucide:code" class="w-6 h-6 text-primary" aria-hidden="true" />
                </div>
                <div>
                  <h3 class="text-xl font-bold mb-2">Development Time</h3>
                  <p class="text-muted-foreground leading-relaxed">
                    Support ongoing development of new features, bug fixes, and improvements to make Playability even better.
                  </p>
                </div>
              </div>
            </Card>

            <!-- Reason 3 -->
            <Card class="dark p-6 lg:p-8 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:border-primary/50">
              <div class="flex items-start gap-4">
                <div class="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0">
                  <Icon name="lucide:database" class="w-6 h-6 text-primary" aria-hidden="true" />
                </div>
                <div>
                  <h3 class="text-xl font-bold mb-2">API & Data Costs</h3>
                  <p class="text-muted-foreground leading-relaxed">
                    Cover expenses for game databases, APIs, and services that power our comprehensive accessibility information.
                  </p>
                </div>
              </div>
            </Card>

            <!-- Reason 4 -->
            <Card class="dark p-6 lg:p-8 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:border-primary/50">
              <div class="flex items-start gap-4">
                <div class="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0">
                  <Icon name="lucide:heart-handshake" class="w-6 h-6 text-primary" aria-hidden="true" />
                </div>
                <div>
                  <h3 class="text-xl font-bold mb-2">Community Impact</h3>
                  <p class="text-muted-foreground leading-relaxed">
                    Help us expand our reach and make gaming more accessible for people with disabilities worldwide.
                  </p>
                </div>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </section>

    <!-- Donation Methods -->
    <section class="py-16 lg:py-20 dark bg-card/30 border-y border-border" aria-labelledby="donate-methods-heading">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto space-y-10">
          <div class="text-center space-y-4 animate-slide-up">
            <h2 id="donate-methods-heading" class="text-3xl lg:text-4xl font-bold">Ways to Support</h2>
            <p class="text-lg text-muted-foreground">
              Choose the method that works best for you
            </p>
          </div>

          <div class="space-y-6 animate-slide-up" style="animation-delay: 150ms">
            <!-- Ko-fi -->
            <Card class="dark p-8 lg:p-10 transition-all duration-300 hover:shadow-2xl hover:border-primary/50 group">
              <div class="flex flex-col md:flex-row items-center gap-8">
                <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-red-500 to-red-600 flex items-center justify-center flex-shrink-0 transition-transform duration-300 group-hover:scale-110">
                  <Icon name="lucide:coffee" class="w-10 h-10 text-white" aria-hidden="true" />
                </div>
                <div class="flex-1 text-center md:text-left space-y-3">
                  <h3 class="text-2xl font-bold">Support on Ko-fi</h3>
                  <p class="text-muted-foreground leading-relaxed">
                    Make a one-time donation or become a monthly supporter. Ko-fi is quick, easy, and supports creators with minimal fees.
                  </p>
                </div>
                <a
                  href="https://ko-fi.com/playability"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center justify-center gap-2 px-8 py-4 bg-gradient-to-r from-red-500 to-red-600 text-white rounded-lg font-semibold transition-all duration-200 hover:scale-105 hover:shadow-lg active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500 focus-visible:ring-offset-2 flex-shrink-0"
                  aria-label="Donate on Ko-fi (opens in new tab)"
                >
                  <Icon name="lucide:heart" class="w-5 h-5" aria-hidden="true" />
                  Donate on Ko-fi
                </a>
              </div>
            </Card>

            <!-- GitHub Sponsors -->
            <Card class="dark p-8 lg:p-10 transition-all duration-300 hover:shadow-2xl hover:border-primary/50 group">
              <div class="flex flex-col md:flex-row items-center gap-8">
                <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center flex-shrink-0 transition-transform duration-300 group-hover:scale-110">
                  <Icon name="lucide:github" class="w-10 h-10 text-white" aria-hidden="true" />
                </div>
                <div class="flex-1 text-center md:text-left space-y-3">
                  <h3 class="text-2xl font-bold">GitHub Sponsors</h3>
                  <p class="text-muted-foreground leading-relaxed">
                    Sponsor the project directly on GitHub. Perfect for developers who want to support open-source accessibility tools.
                  </p>
                </div>
                <a
                  href="https://github.com/sponsors/safegergis"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center justify-center gap-2 px-8 py-4 bg-gradient-to-r from-purple-500 to-purple-600 text-white rounded-lg font-semibold transition-all duration-200 hover:scale-105 hover:shadow-lg active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:ring-offset-2 flex-shrink-0"
                  aria-label="Sponsor on GitHub (opens in new tab)"
                >
                  <Icon name="lucide:heart" class="w-5 h-5" aria-hidden="true" />
                  Sponsor on GitHub
                </a>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </section>

    <!-- Other Ways to Help -->
    <section class="py-16 lg:py-20 dark" aria-labelledby="other-ways-heading">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto space-y-10">
          <div class="text-center space-y-4 animate-slide-up">
            <h2 id="other-ways-heading" class="text-3xl lg:text-4xl font-bold">Other Ways to Help</h2>
            <p class="text-lg text-muted-foreground">
              Not ready to donate? There are still many ways you can support us!
            </p>
          </div>

          <div class="grid md:grid-cols-3 gap-6 animate-slide-up" style="animation-delay: 150ms">
            <!-- Share -->
            <Card class="dark p-6 text-center transition-all duration-300 hover:-translate-y-2 hover:shadow-xl hover:border-primary/50 group">
              <div class="w-16 h-16 rounded-xl bg-primary/10 flex items-center justify-center mx-auto mb-4 transition-transform duration-300 group-hover:scale-110">
                <Icon name="lucide:share-2" class="w-8 h-8 text-primary" aria-hidden="true" />
              </div>
              <h3 class="text-xl font-bold mb-3">Share Playability</h3>
              <p class="text-muted-foreground leading-relaxed">
                Tell others about us on social media and gaming communities
              </p>
            </Card>

            <!-- Contribute -->
            <Card class="dark p-6 text-center transition-all duration-300 hover:-translate-y-2 hover:shadow-xl hover:border-primary/50 group">
              <div class="w-16 h-16 rounded-xl bg-primary/10 flex items-center justify-center mx-auto mb-4 transition-transform duration-300 group-hover:scale-110">
                <Icon name="lucide:git-pull-request" class="w-8 h-8 text-primary" aria-hidden="true" />
              </div>
              <h3 class="text-xl font-bold mb-3">Contribute Code</h3>
              <p class="text-muted-foreground leading-relaxed">
                Help improve Playability by contributing to our open-source codebase
              </p>
            </Card>

            <!-- Submit Reports -->
            <Card class="dark p-6 text-center transition-all duration-300 hover:-translate-y-2 hover:shadow-xl hover:border-primary/50 group">
              <div class="w-16 h-16 rounded-xl bg-primary/10 flex items-center justify-center mx-auto mb-4 transition-transform duration-300 group-hover:scale-110">
                <Icon name="lucide:file-text" class="w-8 h-8 text-primary" aria-hidden="true" />
              </div>
              <h3 class="text-xl font-bold mb-3">Submit Reports</h3>
              <p class="text-muted-foreground leading-relaxed">
                Share your accessibility experiences to help other gamers
              </p>
            </Card>
          </div>
        </div>
      </div>
    </section>

    <!-- Thank You Section -->
    <section class="py-16 lg:py-20 dark bg-gradient-to-t from-background to-neutral-900">
      <div class="container mx-auto px-4">
        <Card class="dark max-w-4xl mx-auto p-8 lg:p-12 text-center animate-fade-in">
          <Icon name="lucide:sparkles" class="w-16 h-16 mx-auto text-primary mb-6" aria-hidden="true" />
          <h2 class="text-3xl lg:text-4xl font-bold mb-4">Thank You!</h2>
          <p class="text-xl text-muted-foreground leading-relaxed max-w-2xl mx-auto">
            Your support, whether financial or otherwise, helps make gaming more accessible for everyone. Together, we're building a more inclusive gaming community.
          </p>
        </Card>
      </div>
    </section>
  </main>
</template>

<script lang="ts" setup>
// Set page metadata
useHead({
  title: 'Support Playability - Donate to Keep Gaming Accessible',
  meta: [
    {
      name: 'description',
      content: 'Support Playability and help us keep gaming accessible for everyone. Your donations help cover server costs, development, and expand our reach to help more gamers with disabilities.'
    }
  ]
});
</script>

<style></style>