RSpec.describe 'Initialize new Gogot repository' do
  let!(:command_result) { `#{command}` }

  context 'without a path' do
    let(:command) { %(gogot init) }

    it 'prints usage information' do
      expect(command_result).to include(('Usage: gogot init <path>'))
    end
  end

  context 'with missing directory' do
    let(:command) { %(gogot init kewl-projekt) }

    before { system('rm -rf ./kewl-projekt') }

    it 'creates the directory if missing' do
      expect(command_result).to include('Initalizing new Gogot repo')
      expect(command_result).to include('Gogot repo initialized in kewl-projekt/.gogot')
    end
  end
end
