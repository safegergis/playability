<template>
  <div class="container mx-auto p-8 py-8">
    <h1 class="text-3xl font-bold mb-6">About Playability</h1>

    <div class="mb-4">
      <p class="text-lg mb-4">
        Playability is a platform dedicated to making gaming accessible to
        everyone. We provide information about game accessibility features and
        help gamers find titles that suit their needs.
      </p>
      <p class="text-lg">
        Our mission is to create a more inclusive gaming community by bridging
        the gap between game developers and players with diverse abilities.
      </p>
    </div>
    <Card class="dark">
      <CardHeader>
        <CardTitle>Frequently Asked Questions</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Accordion type="single" collapsible class="w-full">
            <AccordionItem value="item-1">
              <AccordionTrigger>What is Playability?</AccordionTrigger>
              <AccordionContent>
                Playability is a platform dedicated to making gaming accessible
                to everyone. We provide information about game accessibility
                features and help gamers find titles that suit their needs.
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-2">
              <AccordionTrigger>
                How do you gather information about game accessibility?
              </AccordionTrigger>
              <AccordionContent>
                We collect our basic information from the PC Gaming Wiki. Many
                thanks to them for cataloging thousands of games and features.
                We then allow users to submit their own experiences and updates
                to keep the information as up to date as possible.
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-3">
              <AccordionTrigger>
                Can I contribute to Playability?
              </AccordionTrigger>
              <AccordionContent>
                Yes! We welcome contributions from the community. You can submit
                accessibility information for games, write reviews, or even
                volunteer to help test games for accessibility features.
              </AccordionContent>
            </AccordionItem>
          </Accordion>
          <Accordion type="single" collapsible class="w-full">
            <AccordionItem value="item-1">
              <AccordionTrigger>Who runs Playability?</AccordionTrigger>
              <AccordionContent>
                Playability is built and maintained by a single developer, me,
                with the help of the community. I'm a gamer who wants to make
                gaming more accessible and inclusive. I'm also a computer
                science student who wanted to build something useful for the
                community. build something useful for the community.
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-2">
              <AccordionTrigger>Is Playability free?</AccordionTrigger>
              <AccordionContent>
                Yes! Playability is a free service. We do not collect any
                personal information from our users. If you'd like to support
                us, you can donate to us
                <NuxtLink
                  to="/donate"
                  class="text-blue-100 hover:text-blue-200"
                  target="_blank"
                >
                  here.
                </NuxtLink>
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-3">
              <AccordionTrigger>Is Playability open source?</AccordionTrigger>
              <AccordionContent>
                Yes! You can find the source code for Playability on
                <a
                  href="https://github.com/safegergis/playability"
                  class="text-blue-100 hover:text-blue-200"
                  target="_blank"
                >
                  GitHub.
                </a>
                Feel free to contribute to the project.
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script lang="ts" setup></script>

<style></style>
