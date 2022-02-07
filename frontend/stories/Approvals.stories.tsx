import React, { ComponentProps } from 'react';
import { Approvals } from '@app/Approvals/Approvals';
import { Story } from '@storybook/react';

//👇 This default export determines where your story goes in the story list
export default {
  title: 'Components/Approvals',
  component: Approvals,
};

//👇 We create a “template” of how args map to rendering
const Template: Story<ComponentProps<typeof Approvals>> = (args) => <Approvals {...args} />;

export const ApprovalsStory = Template.bind({});
ApprovalsStory.args = {
  /*👇 The args you need here will depend on your component */
};
