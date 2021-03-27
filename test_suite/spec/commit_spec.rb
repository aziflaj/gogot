RSpec.describe 'Committing changes' do
  around do |example|
    `gogot init kewl-projekt`
    Dir.chdir('kewl-projekt') do
      example.run
    end
    `rm -rf kewl-projekt`
  end

  context 'without a message' do
    let(:command) { `gogot commit` }

    it 'prints usage information' do
      expect(command).to include('Usage: gogot commit <message>')
    end
  end

  context 'with one file in stage' do
    let(:message) { 'Making a commit' }
    let(:command) { `gogot commit #{message}` }
    let(:filename) { 'hello.txt' }

    before do
      File.open(filename, 'w+') { |f| f.write("Howdy y'all") }
      `gogot add .`

      command
    end

    it 'clears index' do
      expect(IO.readlines('.gogot/index').count).to be_zero
    end

    it 'stores a blob for the file' do
      last_commit_id = IO.readlines('.gogot/refs/heads/main').last
      hash_head, hash_tail = last_commit_id[0..1], last_commit_id[2..].tr("\n", '')
      expect(Dir.entries(".gogot/objects/#{hash_head}")).to include(hash_tail)

      tree_hash, _author, _, actual_msg = IO.readlines(".gogot/objects/#{hash_head}/#{hash_tail}")
      expect(message).to eq(actual_msg.tr("\n", ''))

      hash = tree_hash.split.last
      tree_hash_head, tree_hash_tail = hash[0..1], hash[2..].tr("\n", '')
      content = File.read(".gogot/objects/#{tree_hash_head}/#{tree_hash_tail}")

      type, blob_hash, actual_filename = content.split

      expect(type).to eq('blob')
      expect(filename).to eq(actual_filename)

      blob_hash_head, blob_hash_tail = blob_hash[0..1], blob_hash[2..]

      expect(Dir.entries(".gogot/objects/#{blob_hash_head}")).to include(blob_hash_tail)
    end
  end
end
