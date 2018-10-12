import { mount, createLocalVue } from '@vue/test-utils'
import ElementUI from 'element-ui'
import SignupForm from '@/components/SignupForm.vue'

const localVue = createLocalVue()
localVue.use(ElementUI)

describe('SignupForm.vue', () => {
  it('should go top and show message when signup success', () => {
    const signupWrapper = mount(SignupForm, { localVue })
    const button = signupWrapper.find('.button-submit')
    button.trigger('click')
  })
})
