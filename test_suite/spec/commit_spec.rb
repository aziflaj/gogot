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

  context 'with directory in stage' do
    let(:command) { `gogot commit Making a commit` }
    let(:filename) { 'hello.txt' }

    before do
      File.open('file1.txt', 'w+') { |f| f.write('Test File 1') }

      Dir.mkdir('testdir')
      File.open('testdir/file11.txt', 'w+') { |f| f.write('Test File 11 in testdir') }
      File.open('testdir/file12.txt', 'w+') { |f| f.write('Test File 12 in testdir') }
      File.open('testdir/file13.txt', 'w+') { |f| f.write('Test File 13 in testdir') }

      Dir.mkdir('testdir/tester-dir')
      File.open('testdir/tester-dir/file11.txt', 'w+') { |f| f.write('Test File 11 in tester-dir') }
      File.open('testdir/tester-dir/file12.txt', 'w+') { |f| f.write('Test File 12 in tester-dir') }
      File.open('testdir/tester-dir/file13.txt', 'w+') { |f| f.write('Test File 13 in tester-dir') }

      Dir.mkdir('testdir2')
      File.open('testdir2/file21.txt', 'w+') { |f| f.write('Test File 21 in testdir2') }
      File.open('testdir2/file22.txt', 'w+') { |f| f.write('Test File 22 in testdir2') }
      File.open('testdir2/file23.txt', 'w+') { |f| f.write('Test File 23 in testdir2') }

      `gogot add .`

      command
    end

    it 'clears index' do
      expect(IO.readlines('.gogot/index').count).to be_zero
    end

    it 'stores blobs for all files in the directory tree' do
      last_commit_id = IO.readlines('.gogot/refs/heads/main').last
      hash_head, hash_tail = last_commit_id[0..1], last_commit_id[2..].tr("\n", '')
      tree_hash, _author, _, _msg = IO.readlines(".gogot/objects/#{hash_head}/#{hash_tail}")

      tree_hash_head, tree_hash_tail = tree_hash.split.last[0..1], tree_hash.split.last[2..].tr("\n", '')
      content = File.read(".gogot/objects/#{tree_hash_head}/#{tree_hash_tail}")

      content_cells = content.split("\n").map(&:split)

      content_cells.each do |type, hash, name|
        next if type == 'blob'

        hash_head, hash_tail = hash[0..1], hash[2..]
        expect(Dir.entries(".gogot/objects/#{hash_head}")).to include(hash_tail)
        dir_content = File.read(".gogot/objects/#{hash_head}/#{hash_tail}")

        case name
        when 'testdir'
          type, subtree_hash, dirname = dir_content.split("\n").map(&:split).find { |line| line.include?('tree') }
          expect(type).to eq('tree')
          expect(dirname).to eq('tester-dir')

          subtree_hash_head, subtree_hash_tail = subtree_hash[0..1], subtree_hash[2..]
          expect(Dir.entries(".gogot/objects/#{subtree_hash_head}")).to include(subtree_hash_tail)
          subdir_content = File.read(".gogot/objects/#{subtree_hash_head}/#{subtree_hash_tail}")
          expect(subdir_content).to include('file11.txt')
          expect(subdir_content).to include('file12.txt')
          expect(subdir_content).to include('file13.txt')
        when 'testdir2'
          expect(dir_content).to include('file21.txt')
          expect(dir_content).to include('file22.txt')
          expect(dir_content).to include('file23.txt')
        else
          raise "Type unknown #{type}"
        end
      end
    end
  end
end
