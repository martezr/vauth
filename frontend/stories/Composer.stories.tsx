import React, { ComponentProps } from 'react';
import { VirtualMachines } from '@app/VirtualMachines/VirtualMachines';
import { Story } from '@storybook/react';

//👇 This default export determines where your story goes in the story list
export default {
  title: 'Components/VirtualMachines',
  component: VirtualMachines,
};

//👇 We create a “template” of how args map to rendering
const Template: Story<ComponentProps<typeof VirtualMachines>> = (args) => <VirtualMachines {...args} />;

export const ComposerStory = Template.bind({});
ComposerStory.args = {
  /*👇 The args you need here will depend on your component */
};
