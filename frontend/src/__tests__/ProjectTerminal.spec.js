import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import ProjectTerminal from '../components/ProjectTerminal.vue';

describe('ProjectTerminal.vue', () => {
  const sampleGitStatus = {
    branch: 'feat/workspace-manager',
    is_clean: false,
    commits: [
      { hash: 'a1b2c3d', author: 'John Doe', date: '2 hours ago', message: 'feat: add git status' }
    ]
  };

  it('renders active branch and clean/dirty status correctly', () => {
    const wrapper = mount(ProjectTerminal, {
      props: {
        gitStatus: sampleGitStatus,
        consoleLines: [],
        actionRunning: false
      }
    });

    expect(wrapper.text()).toContain('feat/workspace-manager');
    expect(wrapper.text()).toContain('DIRTY');
    expect(wrapper.text()).toContain('a1b2c3d');
    expect(wrapper.text()).toContain('John Doe');
    expect(wrapper.text()).toContain('feat: add git status');
  });

  it('emits action event when buttons are clicked', async () => {
    const wrapper = mount(ProjectTerminal, {
      props: {
        gitStatus: sampleGitStatus,
        consoleLines: [],
        actionRunning: false
      }
    });

    const pullBtn = wrapper.find('button.btn-primary');
    await pullBtn.trigger('click');
    expect(wrapper.emitted().action[0]).toEqual(['pull']);
  });

  it('correctly parses standard ANSI escape sequences into styled classes', () => {
    const wrapper = mount(ProjectTerminal, {
      props: {
        gitStatus: sampleGitStatus,
        consoleLines: [],
        actionRunning: false
      }
    });

    const parseAnsi = wrapper.vm.parseAnsi;
    
    // Test green text: \u001b[32mBuild Success\u001b[0m
    const greenResult = parseAnsi('\u001b[32mBuild Success\u001b[0m');
    expect(greenResult).toHaveLength(1);
    expect(greenResult[0].text).toBe('Build Success');
    expect(greenResult[0].class).toBe('ansi-green');

    // Test red error: \u001b[31mCompilation Failed
    const redResult = parseAnsi('\u001b[31mCompilation Failed');
    expect(redResult).toHaveLength(1);
    expect(redResult[0].text).toBe('Compilation Failed');
    expect(redResult[0].class).toBe('ansi-red');
  });
});
