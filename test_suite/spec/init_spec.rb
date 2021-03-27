RSpec.describe 'Initializing a new Gogot repository' do
  let!(:command_result) { `#{command}` }

  context 'without a path' do
    let(:command) { %(gogot init) }

    it 'prints usage information' do
      expect(command_result).to include('Usage: gogot init <path>')
    end
  end

  context 'with a new directory' do
    let(:command) { %(gogot init kewl-projekt) }

    around do |example|
      `rm -rf ./kewl-projekt`
      example.run
      `rm -rf ./kewl-projekt` # to clean up
    end

    it 'creates the directory if missing' do
      expect(command_result).to include('Initalizing new Gogot repo')
      expect(command_result).to include('Gogot repo initialized in kewl-projekt/.gogot')
    end
  end

  context 'with an empty directory' do
    let(:command) { %(gogot init projekt) }

    around do |example|
      `mkdir projekt`
      example.run
      `rm -rf projekt`
    end

    it 'initializes the repo anyway' do
      expect(command_result).to include('Initalizing new Gogot repo')
      expect(command_result).to include('Gogot repo initialized in projekt/.gogot')
    end
  end
end
