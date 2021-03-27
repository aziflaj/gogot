RSpec.describe 'Logging commits' do
  let(:command_result) { `gogot log` }

  around do |example|
    `gogot init kewl-projekt`
    Dir.chdir('kewl-projekt') do
      example.run
    end
    `rm -rf kewl-projekt`
  end

  context 'with new repo' do
    it 'prints empty log' do
      expect(command_result).to include('Nothing to log on branch')
    end
  end

  context 'with existing repo' do
    before do
      # Commit 1
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
      `gogot commit Commit 1`

      # Commit 2
      File.open('file1.txt', 'a') { |f| f.write('Hello') }
      File.open('testdir/file11.txt', 'a') { |f| f.write('Hello') }

      `gogot add .`
      `gogot commit Commit 2`
    end

    it 'logs commits with hash, author and commit message' do
      command_result.split("\n")[1..].each_with_index do |line, idx|
        _hash, author, msg = line.split('    ')
        expect(msg).to eq("Commit #{idx + 1}")
        expect(author).to include('author')
      end

      printed_msg = command_result.split("\n")[1].split('    ').last
      expect(printed_msg).to eq('Commit 1')

      printed_msg = command_result.split("\n")[2].split('    ').last
      expect(printed_msg).to eq('Commit 2')
    end
  end
end
