import { createLazyFileRoute } from '@tanstack/react-router'

import CTA from '@/components/landing/cta'
import Features from '@/components/landing/features'
import Footer from '@/components/landing/footer'
import Header from '@/components/landing/header'
import Hero from '@/components/landing/hero'
import Pricing from '@/components/landing/pricing'
import Testimonials from '@/components/landing/testimonials'

export const Route = createLazyFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <div className="min-h-screen bg-background font-sans">
      <Header />
      <main>
        <Hero />
        <Features />
        <Testimonials />
        <Pricing />
        <CTA />
      </main>
      <Footer />
    </div>
  )
}
