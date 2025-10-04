<template>
  <main class="min-h-screen">
    <!-- Hero Section -->
    <section class="relative dark overflow-hidden bg-gradient-to-b from-background via-background to-neutral-900 py-16 lg:py-24">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto text-center space-y-6 animate-fade-in">
          <h1 class="text-4xl lg:text-6xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text">
            About Playability
          </h1>
          <p class="text-xl lg:text-2xl text-muted-foreground leading-relaxed">
            Making gaming accessible to everyone, one game at a time
          </p>
        </div>
      </div>

    </section>

    <!-- Mission Section -->
    <section class="py-16 lg:py-20" aria-labelledby="mission-heading">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto space-y-8 animate-slide-up">
          <div class="text-center space-y-4">
            <Icon name="lucide:heart" class="w-12 h-12 mx-auto text-primary animate-pulse" aria-hidden="true" />
            <h2 id="mission-heading" class="text-3xl lg:text-4xl font-bold">Our Mission</h2>
          </div>

          <div class="space-y-6">
            <Card class="dark p-6 lg:p-8 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl">
              <p class="text-lg text-muted-foreground leading-relaxed">
                Playability is a platform dedicated to making gaming accessible to everyone. We provide comprehensive information about game accessibility features and help gamers find titles that suit their needs.
              </p>
            </Card>

            <Card class="dark p-6 lg:p-8 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl">
              <p class="text-lg text-muted-foreground leading-relaxed">
                Our mission is to create a more inclusive gaming community by bridging the gap between game developers and players with diverse abilities. We believe everyone deserves to experience the joy of gaming.
              </p>
            </Card>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-16 dark lg:py-20 bg-card/30 border-y border-border" aria-labelledby="features-heading">
      <div class="container mx-auto px-4">
        <div class="max-w-6xl mx-auto">
          <h2 id="features-heading" class="text-3xl lg:text-4xl font-bold text-center mb-12 animate-slide-up">
            What We Offer
          </h2>

          <div class="grid md:grid-cols-3 gap-8">
            <!-- Feature 1 -->
            <div class="group space-y-4 p-6 rounded-xl bg-card border border-border transition-all duration-300 hover:-translate-y-2 hover:shadow-xl hover:border-primary/50 animate-slide-up">
              <div class="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center transition-transform duration-300 group-hover:scale-110">
                <Icon name="lucide:database" class="w-7 h-7 text-primary" aria-hidden="true" />
              </div>
              <h3 class="text-xl font-bold">Comprehensive Database</h3>
              <p class="text-muted-foreground leading-relaxed">
                Access detailed accessibility information for thousands of games, sourced from PC Gaming Wiki and community reports.
              </p>
            </div>

            <!-- Feature 2 -->
            <div class="group space-y-4 p-6 rounded-xl bg-card border border-border transition-all duration-300 hover:-translate-y-2 hover:shadow-xl hover:border-primary/50 animate-slide-up" style="animation-delay: 100ms">
              <div class="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center transition-transform duration-300 group-hover:scale-110">
                <Icon name="lucide:users" class="w-7 h-7 text-primary" aria-hidden="true" />
              </div>
              <h3 class="text-xl font-bold">Community Driven</h3>
              <p class="text-muted-foreground leading-relaxed">
                Real experiences from real players. Submit reports and help others discover accessible games.
              </p>
            </div>

            <!-- Feature 3 -->
            <div class="group space-y-4 p-6 rounded-xl bg-card border border-border transition-all duration-300 hover:-translate-y-2 hover:shadow-xl hover:border-primary/50 animate-slide-up" style="animation-delay: 200ms">
              <div class="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center transition-transform duration-300 group-hover:scale-110">
                <Icon name="lucide:github" class="w-7 h-7 text-primary" aria-hidden="true" />
              </div>
              <h3 class="text-xl font-bold">Open Source</h3>
              <p class="text-muted-foreground leading-relaxed">
                Transparent and collaborative. View our code on GitHub and contribute to making gaming more accessible.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- FAQ Section -->
    <section class="py-16 lg:py-20" aria-labelledby="faq-heading">
      <div class="container mx-auto px-4">
        <div class="max-w-5xl mx-auto">
          <div class="text-center mb-12 animate-slide-up">
            <Icon name="lucide:help-circle" class="w-12 h-12 mx-auto text-primary mb-4" aria-hidden="true" />
            <h2 id="faq-heading" class="text-3xl lg:text-4xl font-bold">Frequently Asked Questions</h2>
          </div>

          <Card class="dark p-6 lg:p-8 animate-slide-up" style="animation-delay: 150ms">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 lg:gap-8">
              <!-- Column 1 -->
              <Accordion type="single" collapsible class="w-full space-y-2">
                <AccordionItem value="item-1" class="border-border">
                  <AccordionTrigger class="hover:text-primary transition-colors">
                    What is Playability?
                  </AccordionTrigger>
                  <AccordionContent class="text-muted-foreground leading-relaxed">
                    Playability is a platform dedicated to making gaming accessible to everyone. We provide information about game accessibility features and help gamers find titles that suit their needs.
                  </AccordionContent>
                </AccordionItem>

                <AccordionItem value="item-2" class="border-border">
                  <AccordionTrigger class="hover:text-primary transition-colors">
                    How do you gather information about game accessibility?
                  </AccordionTrigger>
                  <AccordionContent class="text-muted-foreground leading-relaxed">
                    We collect our basic information from the PC Gaming Wiki. Many thanks to them for cataloging thousands of games and features. We then allow users to submit their own experiences and updates to keep the information as up to date as possible.
                  </AccordionContent>
                </AccordionItem>

                <AccordionItem value="item-3" class="border-border">
                  <AccordionTrigger class="hover:text-primary transition-colors">
                    Can I contribute to Playability?
                  </AccordionTrigger>
                  <AccordionContent class="text-muted-foreground leading-relaxed">
                    Yes! We welcome contributions from the community. You can submit accessibility information for games, write reviews, or even volunteer to help test games for accessibility features.
                  </AccordionContent>
                </AccordionItem>
              </Accordion>

              <!-- Column 2 -->
              <Accordion type="single" collapsible class="w-full space-y-2">
                <AccordionItem value="item-4" class="border-border">
                  <AccordionTrigger class="hover:text-primary transition-colors">
                    Who runs Playability?
                  </AccordionTrigger>
                  <AccordionContent class="text-muted-foreground leading-relaxed">
                    Playability is built and maintained by a single developer with the help of the community. I'm a gamer who wants to make gaming more accessible and inclusive, and a computer science student who wanted to build something useful for the community.
                  </AccordionContent>
                </AccordionItem>

                <AccordionItem value="item-5" class="border-border">
                  <AccordionTrigger class="hover:text-primary transition-colors">
                    Is Playability free?
                  </AccordionTrigger>
                  <AccordionContent class="text-muted-foreground leading-relaxed">
                    Yes! Playability is a free service. We do not collect any personal information from our users. If you'd like to support us, you can
                    <NuxtLink
                      to="/donate"
                      class="text-primary hover:text-primary/80 underline underline-offset-4 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 rounded-sm"
                    >
                      donate here
                    </NuxtLink>.
                  </AccordionContent>
                </AccordionItem>

                <AccordionItem value="item-6" class="border-border">
                  <AccordionTrigger class="hover:text-primary transition-colors">
                    Is Playability open source?
                  </AccordionTrigger>
                  <AccordionContent class="text-muted-foreground leading-relaxed">
                    Yes! You can find the source code for Playability on
                    <a
                      href="https://github.com/safegergis/playability"
                      class="text-primary hover:text-primary/80 underline underline-offset-4 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 rounded-sm"
                      target="_blank"
                      rel="noopener noreferrer"
                      aria-label="Playability GitHub repository (opens in new tab)"
                    >
                      GitHub
                    </a>.
                    Feel free to contribute to the project.
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
            </div>
          </Card>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="dark py-16 lg:py-20 bg-gradient-to-t from-background to-neutral-900 " aria-labelledby="cta-heading">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto text-center space-y-8 animate-fade-in">
          <h2 id="cta-heading" class="text-3xl lg:text-4xl font-bold">
            Ready to Make Gaming More Accessible?
          </h2>
          <p class="text-xl text-muted-foreground">
            Join our community and help create a more inclusive gaming world.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <NuxtLink
              to="/search"
              class="inline-flex items-center justify-center gap-2 px-8 py-4 bg-primary text-primary-foreground rounded-lg font-semibold transition-all duration-200 hover:scale-105 hover:shadow-lg active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
            >
              <Icon name="lucide:search" class="w-5 h-5" aria-hidden="true" />
              Search Games
            </NuxtLink>
            <a
              href="https://github.com/safegergis/playability"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center justify-center gap-2 px-8 py-4 bg-card border border-border text-foreground rounded-lg font-semibold transition-all duration-200 hover:scale-105 hover:shadow-lg hover:border-primary active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
              aria-label="Contribute on GitHub (opens in new tab)"
            >
              <Icon name="lucide:github" class="w-5 h-5" aria-hidden="true" />
              Contribute on GitHub
            </a>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>

<script lang="ts" setup>
// Set page metadata
useHead({
  title: 'About Playability - Making Gaming Accessible to Everyone',
  meta: [
    {
      name: 'description',
      content: 'Learn about Playability\'s mission to make gaming accessible to everyone. Discover our community-driven approach to game accessibility information.'
    }
  ]
});
</script>

<style></style>
